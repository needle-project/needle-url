#!/bin/bash

# Colors
RED="\033[33;31m";
GREEN="\033[33;32m"
YELLOW="\033[38;5;11m"
RESET="\033[0m"

BASEURL="https://adrian.tilita.ro/"
APPURL="http://127.0.0.1:9293/"
AUTH="test:new"

paddash=$(printf '%0.1s' "-"{1..60})
padspace=$(printf '%0.1s' " "{1..60})
padlength=74

function showresult {
    leftstring=$1
    rightstring=$2
    printf '%s' "$leftstring"
    printf '%*.*s' 0 $((padlength - ${#leftstring} - ${#rightstring} )) "$padspace"
    printf  "[ ${rightstring} ]\n"
}

echo "${YELLOW}Running smoke tests${RESET}"
echo "${YELLOW}${paddash}${RESET}"


# Delete first if any legacy left
curl -o /dev/null --silent -u ${AUTH} -H "Content-Type: application/json" -X DELETE ${APPURL}url/check-redirect

# Create new
HTTP_CALL=$(curl -o /dev/null --silent -u ${AUTH} -H "Content-Type: application/json" -X POST -d '{"from_url": "check-redirect","to_url": "http://example.com"}' --write-out '%{http_code}' ${APPURL}url)
IFS='-' read -r -a HTTP_STATUS <<< "$HTTP_CALL"

leftstring="CREATE NEW"
if [ ${HTTP_STATUS} -ne 201 ]
    then
        rightstring="${RED}CRITICAL${RESET}"
        showresult "${leftstring}" "${rightstring}"
        printf "  - Got ${RED}${HTTP_STATUS}${RESET} instead of ${GREEN}201${RESET}\n"
        exit 2
    else
        rightstring="${GREEN}OK${RESET}"
        showresult "${leftstring}" "${rightstring}"
fi

# Create duplicate
HTTP_CALL=$(curl -o /dev/null --silent -u ${AUTH} -H "Content-Type: application/json" -X POST -d '{"from_url": "check-redirect","to_url": "http://example.com"}' --write-out '%{http_code}' ${APPURL}url)
IFS='-' read -r -a HTTP_STATUS <<< "$HTTP_CALL"

leftstring="CREATE DUPLICATE"
if [ ${HTTP_STATUS} -ne 409 ]
    then
        rightstring="${RED}CRITICAL${RESET}"
        showresult "${leftstring}" "${rightstring}"
        printf "  - Got ${RED}${HTTP_STATUS}${RESET} instead of ${GREEN}409${RESET}\n"
        exit 2
    else
        rightstring="${GREEN}OK${RESET}"
        showresult "${leftstring}" "${rightstring}"
fi

# Call
HTTP_CALL=$(curl -o /dev/null --silent -u ${AUTH} -H "Content-Type: application/json" --write-out '%{http_code}' ${APPURL}check-redirect)
IFS='-' read -r -a HTTP_STATUS <<< "$HTTP_CALL"

leftstring="REDIRECT URL"
if [ ${HTTP_STATUS} -ne 301 ]
    then
        rightstring="${RED}CRITICAL${RESET}"
        showresult "${leftstring}" "${rightstring}"
        printf "  - Got ${RED}${HTTP_STATUS}${RESET} instead of ${GREEN}301${RESET}\n"
        exit 2
    else
        rightstring="${GREEN}OK${RESET}"
        showresult "${leftstring}" "${rightstring}"
fi

# Update
HTTP_CALL=$(curl -o /dev/null --silent -u ${AUTH} -H "Content-Type: application/json" -X PATCH -d '{"from_url": "check-redirect","to_url": "https://www.github.com"}' --write-out '%{http_code}' ${APPURL}url)
IFS='-' read -r -a HTTP_STATUS <<< "$HTTP_CALL"

leftstring="UPDATE URL"
if [ ${HTTP_STATUS} -ne 200 ]
    then
        rightstring="${RED}CRITICAL${RESET}"
        showresult "${leftstring}" "${rightstring}"
        printf "  - Got ${RED}${HTTP_STATUS}${RESET} instead of ${GREEN}200${RESET}\n"
        exit 2
    else
        rightstring="${GREEN}OK${RESET}"
        showresult "${leftstring}" "${rightstring}"
fi

# Update
HTTP_CALL=$(curl -o /dev/null --silent -u ${AUTH} -H "Content-Type: application/json" -X DELETE --write-out '%{http_code}' ${APPURL}url/check-redirect)
IFS='-' read -r -a HTTP_STATUS <<< "$HTTP_CALL"

leftstring="DELETE URL"
if [ ${HTTP_STATUS} -ne 200 ]
    then
        rightstring="${RED}CRITICAL${RESET}"
        showresult "${leftstring}" "${rightstring}"
        printf "  - Got ${RED}${HTTP_STATUS}${RESET} instead of ${GREEN}200${RESET}\n"
        exit 2
    else
        rightstring="${GREEN}OK${RESET}"
        showresult "${leftstring}" "${rightstring}"
fi

# Call
HTTP_CALL=$(curl -o /dev/null --silent -u ${AUTH} -H "Content-Type: application/json" --write-out '%{http_code}' ${APPURL}check-redirect)
IFS='-' read -r -a HTTP_STATUS <<< "$HTTP_CALL"

leftstring="REDIRECT URL"
if [ ${HTTP_STATUS} -ne 307 ]
    then
        rightstring="${RED}CRITICAL${RESET}"
        showresult "${leftstring}" "${rightstring}"
        printf "  - Got ${RED}${HTTP_STATUS}${RESET} instead of ${GREEN}307${RESET}\n"
        exit 2
    else
        rightstring="${GREEN}OK${RESET}"
        showresult "${leftstring}" "${rightstring}"
fi