
chrome.tabs.onUpdated.addListener((tabId, changeInfo, tab) => {
    console.log('Tab updated:', { tabId, changeInfo, tab });
    
    if (changeInfo.url) {
      const url = changeInfo.url;
      // Option 1: Simple substring check (case insensitive)
      if (url.toLowerCase().includes('checkout')) {
        console.log(`Tab ${tabId} is a checkout page: ${url}`);
        // You can dispatch a message or perform other actions here
      } else {
        console.log(`Tab ${tabId} is not a checkout page.`);
      }
    }
  });