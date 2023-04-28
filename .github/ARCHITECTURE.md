# Stable Release Strategy

## Introduction

This document describes a strategy to manage the development of software solutions in this repository, specifically focusing on new feature creation, breaking changes, and refactorings. The goal of this strategy is to ensure high-quality software delivery by establishing boundaries, structure, and an agreement process.

To avoid rushed decisions and complex code management, it is important to reach a community consensus before implementing solutions. This strategy emphasizes the documentation of architectural decisions, their rationale, and trade-offs. It also promotes communication, collaboration, traceability, and accountability in the architectural evolution of the project.

By following this strategy, we aim to improve the quality, consistency, and maintainability of the software solutions. Additionally, we hope to foster a culture of transparency, learning, and continuous improvement among the project developers.

This document solely focuses on the strategy and architecture discussions and decisions and does not cover the repository structure and code.

Legend:
- `ADR` - Architecture Decision Record - a document that captures an important architectural decision made along with its context and consequences.
- `ADL` - Architecture Decision Log - a collection of `ADR`s persisted in this repository as a changelog with structured folders.
- `ADR Template` - a template format that should be followed in all `ADR` documents. It can be found at [ADR_TEMPLATE.md](ADR_TEMPLATE.md) file.
- `POC` - proof of concept - mainly code that exposes API or a naive implementation that is not expected to work completely correctly, but to visualize the concept.

## Structure

- Issues, PRs, and Discussions - all of them are equal sources of information and are related, meaning the root topic should contain all relationships between conversations related to the topic from other sources. These should be used to express a need in codebase change or start discussion for improvements in project.

- `./{package}/docs` - folder for package documentation and `ADL`.
  - `/{RRR-MM}` - folder with multiple sub-directories containing `ADR`'s.
    - `{title}` - folder with specific `ADR` contents like images, diagrams and `ADR` document(README.md).
      - `/README.md` - the record itself following the `ADR Template`.

## Architecture Decision Processing

The process of creating and maintaining ADL and ADR involves the following steps:

- **Identify an architectural decision** that needs to be made or documented.
- **Create** own **branch** with `adr/` prefix
- **Create** a new **ADR** using a predefined template that includes the following sections: title, status, context, decision, consequences, alternatives considered, and references.
- **Create directory** with a current date `ADR` directory if it do not exist. *(e.g 2023-01)*
- **Create directory** with `ADR` title *(e.g memory-cache-algorithm)*
- **Add the ADL contents** in created titled directory like `README.MD` *(actual ADR file)* and any images.
- **Push Pull Request** with ADR contents. The proof of concept code is recommended but **it should not be on the same branch with ADR**.
- **Review and update the ADR** content, description and concerns as needed, changing its status from proposed to accepted, rejected, deprecated, or superseded.
- Communicate and share the ADR with the relevant parties and **solicit feedback and approval**.
- If you lack of time or you're not interested in making the codebase changes yourself - open a question to get a candidate. It's mandatory to have a volunteer before merging pull request *(you can assign `good first issue` label to PR for such scenario)*.
- **Merge Pull Request** and assign related person to the issue (or create one) that will mean the content is ready to be implemented.

Pull requests based on the `/adr` branch are used to propose and discuss changes to the Architecture Decision Records (ADRs), while other pull requests may be used to propose and discuss changes to the codebase, **ADR proof of concept code** or other aspects of the repository. You can read about the branching convention in [CONTRIBUTING.md](CONTRIBUTING.md#branch-convention).

An important aspect of managing Architecture Decision Records (ADRs) is ensuring that they remain specific and immutable. Each ADR should focus on a single architectural decision, rather than addressing multiple decisions. Additionally, once an ADR has been merged, its existing information should not be altered. Instead of modifying an existing ADR, a new ADR can be created to supersede it. This approach helps to maintain the integrity and traceability of the decision-making process.

## Sample Timeline

Here is an example timeline illustrating the entire process, using a memory cache implementation for the `caches` package:

1. Open up an issue titled `Memory cache implementation`.
2. Discuss the problem and solution.
3. Create a branch `adr/caches/memory-implementation`.
4. Push directory and README for `ADR` => `./caches/docs/2023-05/memory-cache/README.md`. The README file should be filled up based on `ADR Template`. Initital `ADR` status should be `Proposed`. The content does not need to be completely filled up, but you can propose your solutions and decision.
5. **(optional)** include any `POC` by creating separate branch `poc/caches/memory-implementation` and refer the branch in `adr/` pull request.
6. Create Pull Request in `Draft` state.
7. Proceed with any community questions, concerns and update the ADR with relevant additional information or even change the decision, include another alternatives.
8. Get an approval from relevant party and merge `ADR`.
9. Close POC *(if such pull request was created)*.
10. Create a branch `feat/caches/memory-implementation`.
11. Implement the solution based on decision in `ADR`.
12. Create Pull Request & wait for approve & merge.

Events that may occure in future (e.g better data structure):
1. Open up an issue `Stale memory cache implementation`.
2. Discuss the enhancement on more modern or faster approach.
3. Create a branch `adr/caches/memory-data-structure`.
4. Push directory and README for `ADR` => `./caches/docs/2029-05/memory-data-structure/README.md`. The README file should be filled up based on `ADR Template`. Initital `ADR` status should be `Proposed`. The content does not need to be completely filled up, but you can propose your solutions and decision. This ADR should also refer to `./caches/docs/2023-05/memory-cache` `ADR` where previous data-structure was decided.
5. **(optional)** include any `POC` and preferably benchmark in this case.
- ... (the rest steps (6-12) remain the same)
