# Release Notes Aggregator (rna)
Aggregates release notes based on the the RELEASE.md inside a repository.
Pulls release notes in based on a yaml file that you specify, and then
hosts it locally for you.
## Installing
`go get github.com/lestopher/rna`
## Usage
`rna -port=":8888" -conf="/etc/rna/repos.yml"`
