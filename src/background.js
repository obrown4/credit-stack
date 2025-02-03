chrome.tabs.onUpdated.addListener(async (tabId, changeInfo, tab) => {
    console.log("Tab updated:", tab.url);

    // Ensure the page is fully loaded and avoid restricted pages
    if (changeInfo.status === 'complete') {
        const result = await chrome.tabs.sendMessage(tabId, { type: "NEW_PAGE" });
    }
});
