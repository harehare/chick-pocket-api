run:
	go run main.go pocket.go limit.go

dist:
	git push heroku master

deps:
	dep ensure
