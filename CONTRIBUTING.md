# How to contribute

We're happy you're considering contributing to HTMLhouse!

It won't take long to get up to speed on everything. These are our development resources:

* We keep track of tasks in [Phabricator](https://phabricator.write.as/tag/htmlhouse/). Feel free to [sign up](https://phabricator.write.as/auth/start/?next=%2Ftag%2Fhtmlhouse%2F).
* We accept and respond to bugs here on [GitHub](https://github.com/writeas/htmlhouse/issues).
* We're usually in #writeas on freenode, but if not, find us on our [Slack channel](http://slack.write.as).

## Submitting changes

Please send a [pull request](https://github.com/writeas/htmlhouse/compare) off the **develop** branch with a clear explanation of what you've done.

Please follow our coding conventions below and make sure all of your commits are atomic. Larger changes should have commits with more detailed information on what changed, any impact on existing code, rationales, etc.

## Coding conventions

We strive for consistency above all. Reading the small codebase should give you a good idea of the conventions we follow.

### Go

* We use `go fmt` before committing anything
* We aim to document all exported entities
* Go files are broken up into logical functional components
* General functions are extracted into modules when possible

### Javascript

* We use tabs?!?
