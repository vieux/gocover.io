[GoCover.io](http://gocover.io)
===============================

Shutting down...
----------------

The GoCover project started in 2014 and shutdown in 2023.

What was it?
------------

[GoCover.io](http://gocover.io) was a service that offered code coverage of any golang package.

How did it work?
----------------

GoCover gathered code coverage by executing a package's tests within an isolated [Docker](http://docker.io) container.

What's Next?
------------

[@ncruces](https://github.com/ncruces) released `go-coverage-report`, a [GitHub Action](https://github.com/marketplace/actions/go-coverage-report) that will generate goverage report on every commit and store it in Github (hidden in the your project's wiki).


You add this to your CI pipeline (_after_ tests are run):
```yaml
    - name: Update coverage report
      uses: ncruces/go-coverage-report@main
```

Then the badge is applied to the `README.md` as such:
```markdown
[![Go Coverage](https://github.com/USER/REPO/wiki/coverage.svg)](https://raw.githack.com/wiki/USER/REPO/coverage.html)

You can see it in action on: [`github.com/ncruces/julianday`](https://github.com/ncruces/julianday)