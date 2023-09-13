const iframeConfigName = "iframeConfig";
async function getConfig(file) {
  const response = await fetch(file);
  const configTemp = await response.json();
  const IFrameMap = new Map();
  let i = 0;
  while (i < configTemp.Iframes.length) {
    console.log(configTemp.Iframes[i]);
    IFrameMap.set(configTemp.Iframes[i].frameName, configTemp.Iframes[i]);
    i++;
  }
  const mapArray = Array.from(IFrameMap);
  const jsonObject = Object.fromEntries(mapArray);
  const jsonConfig = JSON.stringify(jsonObject);
  localStorage.setItem(iframeConfigName, jsonConfig);
}
