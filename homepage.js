const name = document.getElementById("topicName")
const askButton = document.getElementById("ask-button")

askButton.addEventListener("click", function () {
    fetch("http://localhost:8000/api/publish", {
        method: "POST",
        body: JSON.stringify({
            Name: 'cherry',
            Topic: document.getElementById("topicName").value,
            Msg: document.getElementById("name").value,
        }),
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            console.log(result)
        });
    }).catch((error) => {
        console.log(error)
    });
})

