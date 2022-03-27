document.addEventListener('DOMContentLoaded', function () {
        var table = document.getElementById("mytable")
        const topicnamesearch = document.getElementById("topicnamesearch")
        const showTopicSearchData = document.getElementById("showtopicsearchdata")
        const searchButton = document.getElementById("topicsearchbutton")
        searchButton.addEventListener('click', function() {
            console.log("this is working")
            fetch("http://localhost:8000/api/fetch", {
                method: "POST",
                body: JSON.stringify({
                    Topic: document.getElementById("topicnamesearch").value,
                }),
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json'
                },
            }).then((response) => {
                response.text().then(function (data) {
                    let result = JSON.parse(data);
                    console.log(result)

                    // var updatesstr = ""
                    // for(i=0;i<result.Updates.length;i++) {
                    //     updatesstr += "\n"+"device _id = "+result.Updates[i]._id + "device name = "+result.Updates[i].Name + "device status = "+result.Updates[i].Status + "time of update = "+"device name = "+result.Updates[i].TimeStampFE   
                    // }

                    var row = table.insertRow(0)
                    row.appendChild(document.createElement('th')).innerHTML = "SNO"
                    row.appendChild(document.createElement('th')).innerHTML = "NAME"
                    row.appendChild(document.createElement('th')).innerHTML = "STATUS"
                    row.appendChild(document.createElement('th')).innerHTML = "TIMESTAMP OF LAST UPDATE"

                    for(i=0;i<result.Updates.length;i++) {
                        var row = table.insertRow(i+1)
                        row.insertCell(0).innerHTML = result.Updates[i]._id
                        row.insertCell(1).innerHTML = result.Updates[i].Name
                        row.insertCell(2).innerHTML = result.Updates[i].Status
                        row.insertCell(3).innerHTML = result.Updates[i].TimeStampFE
                    }

                    //showTopicSearchData.textContent = "transitions of topic "+ result.Topic + "\n" + updatesstr
                });
            }).catch((error) => {
                console.log(error)
        });
    });
});

