#!/bin/bash
URL='http://localhost:8080/t/demo/crud'

function jsoncurl() {
  curl -v -H "Content-Type: application/json" $@
}

action=`echo $1 | tr '[:upper:]' '[:lower:]'`
case $action in
  "retrieve")
    echo "Doing $1"
    data="{\"id\":$2}"
    if [ "$#" -ne 2 ]; then
      data="{}"
    fi

    jsoncurl -X GET $URL/get -d $data
    ;;
  "create")
    echo "Doing $1"
    if [ "$#" -ne 3 ]; then
      printf "You need to enter name and last name.\nUSAGE: $0 $1 <name> <last_name>\n"
      exit 0
    fi

    data="{\"name\":\"$2\",\"last_name\":\"$3\"}"
    jsoncurl -X POST $URL/create -d $data
    ;;
  "update")
    echo "Doing $1"
    if [ "$#" -ne 4 ]; then
      printf "You need to enter ID, name and last name.\nUSAGE: $0 $1 <id> <name> <last_name>\n"
      exit 0
    fi

    data="{\"id\":$2,\"name\":\"$3\",\"last_name\":\"$4\"}"
    jsoncurl -X PUT $URL/update -d $data
    ;;
    "delete")
      echo "Doing $1"
      if [ "$#" -ne 2 ]; then
        printf "You need to enter ID\n"
        exit 0
      fi

      data="{\"id\":$2}"
      jsoncurl -X DELETE $URL/delete -d $data
      ;;
  *)
    echo "Action not found '$1'"
    echo "Options are create, retrieve, delete, update"
    ;;
esac
