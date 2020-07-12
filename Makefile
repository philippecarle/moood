make init
# Bootstrap your application (fetch some data files, make some API calls, request user input etc...)

make style
# Check lint, code styling rules. e.g. pylint, phpcs, eslint, style (java) etc ...

make complexity
# Cyclomatic complexity check (McCabe), radon (python), eslint (js), PHPMD, rules (scala) etc ...

make format
# Format code. e.g Prettier (js), format (golang)

make test
# Shortcut to launch all the test tasks (unit, functional and integration)

make test-unit
# launch unit tests. e.g. pytest, jest, phpunit, JUnit etc...

make test-functional
# launch functional tests. e.g behat, JBehave, Behave, CucumberJS, Cucumber etc...

make test-integration
# launch integration tests. e.g pytest, jest, phpunit etc...

make security
# Shortcut to launch all the security tasks (security-sast, security-dependency-scan)

make security-sast
# launch static application security testing (SAST). e.g Gosec, bandit, Flawfinder, NodeJSScan, phpcs-security-audit, brakeman.

make security-dependency-scan
# launch a dependency scanning to trigger know vulnerabilities. e.g Retire.js, gemnasium, bundler-audit.

make run
# Locally run the application, e.g. node index.js, python -m myapp, go run myapp etc ...

make watch
# Hot reloading for development.

make build
	docker-compose up -d --build --remove-orphans
