chrome.tabs.onUpdated.addListener(async (tabId, changeInfo, tab) => {
    console.log("Tab updated:", tab.url);

    // Ensure the page is fully loaded and avoid restricted pages
    if (changeInfo.status === 'complete') {
        try {
            const result = await chrome.tabs.sendMessage(tabId, { type: "NEW_PAGE" });

            if (result) {
                console.log("Message sent successfully:", result);
                chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
                    if (request.action === "notifyUser") {
                        chrome.notifications.create({
                            type: 'basic',
                            iconUrl: 'icon.png',  // Use your extension's icon here
                            title: 'Credit Stack',
                            message: request.data,
                            priority: 2
                        });
                    }
                });
            }
        } catch (error) {
            console.error("Error sending message to content script:", error);
        }
    } else {
        console.log("Skipping non-HTTP page or incomplete load.");
    }
});




// // Use MutationObserver to detect changes in dynamic pages
// const observer = new MutationObserver(() => {
//     console.log("Mutation here!");
//     if (isCheckoutPage()) {
//         console.log("Send message!");
//         chrome.tabs.sendMessage(tabId, { type: "CHECKOUT_DETECTED" });
//         // if (response){
//         //     console.log("Message sent successfully");
//         //     reboot();
//         // }
//         observer.disconnect();  // Stop observing after detection to save resources
//     }
// });

// chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
//     console.log("Message received");
//     console.log(request);
//     if (request.type === "CHECKOUT_DETECTED") {
//         console.log("CHECKOUT_DETECTED");


//         const domain = new URL(sender.tab.url).hostname;
//         console.log(domain);
//         getRecommendataion(domain);
//     }
// });


