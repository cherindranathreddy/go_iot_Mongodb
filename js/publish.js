const name = document.getElementById("topicName")
const publishButton = document.getElementById("ask-button")
const showData = document.getElementById("show-data")
const on_off = document.getElementById("switch")

var status = "off"
on_off.addEventListener("click",function() {
    if(status == "off") status="on"
    else status="off"
    console.log(status)
});


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
            showData.textContent = "The device "+ result.Name + " with id="+result.Id + " is turned "+result.Status+"\n We decided to turn it " +result.Status+ " at time="+ result.TimeStampFE + " and it got turned on at time="+ result.TimeStampBE;
            //showData.textContent = "The device "+ result.Name + " is turned "+result.Status+"\n We decided to turn it " +result.Status+ " at time="+ result.TimeStampFE + " and it got turned on at time="+ result.TimeStampBE;
        });
    }).catch((error) => {
        console.log(error)
    });
})

// var table = document.getElementById("mytable")
// const searchButton = document.getElementById("topicsearchbutton")
// const showTopicSearchData = document.getElementById("showtopicsearchdata")
// searchButton.addEventListener("click", function() {
//     fetch("http://localhost:8000/api/fetch", {
//         method: "POST",
//         body: JSON.stringify({
//             Topic: document.getElementById("topicnamesearch").value,
//         }),
//         headers: {
//             'Accept': 'application/json',
//             'Content-Type': 'application/json'
//         },
//     }).then((response) => {
//         response.text().then(function (data) {
//             let result = JSON.parse(data);
//             console.log(result)

//             // var updatesstr = ""
//             // for(i=0;i<result.Updates.length;i++) {
//             //     updatesstr += "\n"+"device _id = "+result.Updates[i]._id + "device name = "+result.Updates[i].Name + "device status = "+result.Updates[i].Status + "time of update = "+"device name = "+result.Updates[i].TimeStampFE   
//             // }

//             var row = table.insertRow(0)
//             row.appendChild(document.createElement('th')).innerHTML = "SNO"
//             row.appendChild(document.createElement('th')).innerHTML = "NAME"
//             row.appendChild(document.createElement('th')).innerHTML = "STATUS"
//             row.appendChild(document.createElement('th')).innerHTML = "TIMESTAMP OF LAST UPDATE"

//             for(i=0;i<result.Updates.length;i++) {
//                 var row = table.insertRow(i+1)
//                 row.insertCell(0).innerHTML = result.Updates[i]._id
//                 row.insertCell(1).innerHTML = result.Updates[i].Name
//                 row.insertCell(2).innerHTML = result.Updates[i].Status
//                 row.insertCell(3).innerHTML = result.Updates[i].TimeStampFE
//             }

//             //showTopicSearchData.textContent = "transitions of topic "+ result.Topic + "\n" + updatesstr
//         });
//     }).catch((error) => {
//         console.log(error)
//     });
// })




