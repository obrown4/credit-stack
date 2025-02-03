const stores = new Map(
    [["www.walmart.com", "groceries"],
    ["www.ubereats.com", "dining"],
    ["www.delta.com", "travel"]]
);

const cards = new Map(
    [["Chase Saphire Preferred", [3, 3, 5]],
    ["Discover It", [1, 5, 1]],
    ["Bilt Card", [2, 1, 1]]
]);

const indexes = new Map(
    [["groceries", 0],
    ["dining", 1],
    ["travel", 2]]
);

chrome.runtime.onMessage.addListener((request, sender, response) => {
    console.log("Message received");
    chrome.storage.local.set({card_name: "Finding best card..."});
    if(request.type == "NEW_PAGE" && isCheckoutPage()){
        console.log("CHECKOUT_DETECTED");
        const url = window.location
        getRecommendataion(url.hostname);
        response({status: "success"});
    }
    response({status: "failure"});
});


function isCheckoutPage() {
    console.log("Checking if it's a checkout page...");
    
     // 1. Check URL patterns first (most reliable)
     const checkoutUrlPatterns = [
        /checkout/i,
        /payment/i,
        /order\/confirm/i,
        /cart.*confirm/i,
        /billing/i,
        /complete.*purchase/i
    ];
    
    const currentUrl = window.location.href.toLowerCase();
    for (const pattern of checkoutUrlPatterns) {
        if (pattern.test(currentUrl)) {
            console.log("Checkout detected via URL pattern");
            return true;
        }
    }

    
    console.log("Not a checkout page");
    return false;
}

function getRecommendataion(domain){
    console.log(domain);
    let recommendation = document.getElementById("recommendation");
    if(!domain in stores){
        console.log("No store found");
        chrome.storage.local.set({card_name: "Any card will do"});
        return;
    }

    let spend_category = stores.get(domain);
    let index = indexes.get(spend_category);

    console.log("Spend category: " + spend_category);

    let max_points = 0;
    let card_name = "";

    for (let [key, value] of cards){
        console.log(key + " " + value);
        console.log("Max points: " + max_points);
        if (value[index] > max_points){
            max_points = value[index];
            card_name = key;

            console.log("Max points: " + max_points);
        }
    }

    console.log("Recommended card: " + card_name);
    chrome.storage.local.set({card_name: card_name});
    chrome.tabs.create({url: 'popup.html'}) 

    
}
