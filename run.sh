# Runs trakx prod mode in screen
# todo: write pid to a file and all that jazz so that I don't need screen

screen -L -dm bash -c "while true; do go run trakx.go; done"
screen -list
