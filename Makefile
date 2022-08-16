swag:
	swag fmt && swag init
	curl -X POST https://converter.swagger.io/api/convert -d @docs/swagger.json --header 'Content-Type: application/json' > docs/swagger3.json
	mv docs/swagger3.json docs/swagger.json

remove_local:
	git remote update --prune
	git checkout origin/dev
	git for-each-ref --format '%(refname:short)' refs/heads | xargs git branch -D