
const testButton = document.getElementById("test-button")
const time = new EventSource('/time');
const buttonTime = document.getElementById("button-time")
const name = document.getElementById("name")
const askButton = document.getElementById("ask-button")
const askTime = document.getElementById("ask-time")



time.addEventListener('time', (e) => {
    document.getElementById("actual-time").innerHTML = "Actual time using SSE: " + e.data;
}, false);



testButton.addEventListener("click", function () {
    console.log("Button UNO clicked")
    buttonTime.textContent = "Button clicked at: " + new Date().toLocaleString("en-IE")
});



askButton.addEventListener("click", function () {
    let data = {
        Name: name.value,
        Time: new Date().toLocaleString("en-IE"),
    };
    fetch("/get_time", {
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            console.log(result);
            askTime.textContent = "Backend processing took " + result["Duration"] //+ " and ended at " + result["Time"]
        });
    }).catch((error) => {
        console.log(error);
    });
})

