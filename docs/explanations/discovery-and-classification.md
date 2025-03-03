---
title: How detection and classification works
---

# Sensitive data flow

Data discovery and data classification are the core of how Bearer CLI identifies sensitive data flows in your codebase.

## How does Bearer CLI discover and classify data?

When you run Bearer CLI on a codebase, it discovers and classifies data by identifying patterns in the source code. Specifically, it looks for data types and matches against them. Most importantly, it never views the actual values—but only the code itself.

Bearer CLI supports 120+ data types from sensitive data categories such as Personal Data (PD), Sensitive PD, Personally identifiable information (PII), Personal Health Information (PHI). You can view the full list in the [supported data types documentation](/reference/datatypes/).

## Data Discovery

Discovery is the first step of the process. Bearer CLI uses static code analysis in two ways.

- Analyzing class **names, methods, functions, variables, properties, and attributes**. It then ties those together to detected data structures. It goes as far as doing variable reconciliation.
- Analyzing **data structure definitions files** such as OpenAPI, SQL, GraphQL, and Protobuf files.

## Data Classification

Bearer CLI next applies pattern matching and heuristics to classify the discovered data. These patterns and heuristics were built by analyzing tens of thousands of pieces of training sets from open-source projects, as well as Bearer's design partners. The heuristics rely heavily on [cluster analysis](https://en.wikipedia.org/wiki/Cluster_analysis) from the retrieved contextual information.

### Example 1: Known object representing an individual

```txt
Object: User

Attributes:
  - id
  - name
  - email
```

For this example, Bearer CLI will check whether the data is linked to a person or not. We call this a **Known Object**. It's an object that we pre-identified as corresponding to an individual. Some common known objects include user, customer, patient, doctor, and many more.

Bearer CLI starts by checking if the detected object matches a known object. In this case, `User` matches the common user type.

Next, it applies a set of **Known Object Patterns** to check if the attributes of the object match any of the supported data types. In this case, `id` is classified as a User Identifier, `name` as a Full Name, and `email` as an Email Address.

### Example 2: Unknown object representing an individual

```txt
Object: Foo

Attributes:
  - id
  - first_name
  - email
```

For this example, the object `Foo` doesn't match any known object so Bearer CLI can't say with confidence that it represents an individual.

Rather than applying the patterns for known objects, it applies a set of **Unknown Object Patterns** to check if the attributes match any data types. This pattern set is more conservative, as it has a lower confidence that `Foo` represents a person. Once these patterns return a positive result, like they would on `first_name`, Bearer CLI will then assume that this object represents an individual. From there, it will apply an extended set of patterns, **Unknown Object Extended Patterns**, to the remaining attributes to see if they match any data types.

### Example 3: Object linked to an individual with known attributes

```txt
Object: Comment

Attributes:
- id
- user_id
- message
```

In this example, Bearer CLI doesn't match `Comment` with any known object and the unknown object patterns don't return a positive result. This means it is confident that this object is not an individual.

From here, it applies a new set of patterns, called **User Identifier Patterns**, to the attributes as a means to determine if `Comment` is linked to an individual. `user_id` returns a positive result which signals to Bearer CLI that the object is linked to an individual.

Bearer CLI will then apply a set of **Associated Object Patterns** to check if the attributes match a data type. As a result, `message` is classified as "Interactions."

### Example 4: Object linked to an individual with unknown attributes

```txt
Object: Appointment

Attributes:
- Id
- patient_id
- date
- bar
- status
```

In this example:

- Bearer CLI doesn't match `Appointment` with any known object.
- The **Unknown Object Patterns** return negative results.
- The **User Identifier Patterns** return positive based on the `patient_id` attribute.
- The **Associated Object Patterns** return negative results.

In other words, Bearer CLI can't match the attributes to a specific data type.

As a last resort, it applies a new set of patterns called **Known Data Object Patterns**. These attempt to match the object directly to data types, using the attributes as context. In this instance, Bearer CLI understands that the `patient_id` is a medical context and therefore classifies the entire object as a "Health record."  
  
*Side note on PHI: You can tell Bearer CLI that your code is processing health data, and it will make sure to flag information such as email or first_name as PHI instead of PD. To do so, set the `--context` flag to `health`.*  
  
Each object will then need to be reviewed and classified manually by users. This approach avoids false negatives for sensitive data categories.

## How accurate is Bearer CLI's data classification?

Bearer CLI is built with hundreds of patterns for each data type and context object to minimize false positives at the potential expense of missing some true positives.  
  
Our latest tests of unknown data sets **yield an accuracy of 93%**, or 7% false positives.

## You can improve Bearer CLI's detection and classification

Bearer CLI is extensible. Just as you can add specific rules that work for your own needs, you can also new data types, detectors, and classification patterns. Documentation on how to do so is **coming soon**, but for now, you can [browse the repo]({{meta.sourcePath}}) to see how these are implemented.
