document.addEventListener("DOMContentLoaded", function () {
    document.getElementById("recommendation").innerHTML = "Loading...";
    chrome.storage.local.get(["card_name"], function (result) {
        makeRecommendation(result.card_name);
    });
});

function makeRecommendation(card_name) {
    document.getElementById("recommendation").innerHTML = card_name;
}