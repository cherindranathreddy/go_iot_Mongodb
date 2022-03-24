const name = document.getElementById("topicName")
const publishButton = document.getElementById("ask-button")
const showData = document.getElementById("ask-data")
const on_off = document.getElementById("switch")

var status = "off"
on_off.addEventListener("click",function(e) {
    if(status == "off") status="on"
    else status="off"
})


publishButton.addEventListener("click", function () {
    fetch("http://localhost:8000/api/publish", {
        method: "POST",
        body: JSON.stringify({
            Name: document.getElementById("devicename").value,
            Status: status,
            Topic: document.getElementById("topicName").value,
            TimeStampFE: String(new Date()),
        }),
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            result.TimeStampBE = String(new Date())
            console.log(result)
            showData.textContent = "The device "+ result.Name + " is turned "+result.Status+"\n We decided to turn it on at time="+ result.TimeStampFE + " and it got turned on at time="+ result.TimeStampBE;
            //showData.textContent = "The device "+ result.Name + " with id="+result.Id + " is turned "+result.Status+"\n We decided to turn it on at time="+ result.TimeStampFE + " and it got turned on at time="+ result.TimeStampBE;
        });
    }).catch((error) => {
        console.log(error)
    });
})

const searchButton = document.getElementById("topicsearchbutton")
const showTopicSearchData = document.getElementById("showtopicsearchdata")
searchButton.addEventListener("click", function() {
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

            var updatesstr = ""
            for(i=0;i<result.Updates.length;i++) {
                updatesstr += "\n"+"device _id = "+result.Updates[i]._id + "device name = "+result.Updates[i].Name + "device status = "+result.Updates[i].Status + "time of update = "+"device name = "+result.Updates[i].TimeStampFE   
            }

            showTopicSearchData.textContent = "transitions of topic "+ result.Topic + "\n" + updatesstr
        });
    }).catch((error) => {
        console.log(error)
    });
})




