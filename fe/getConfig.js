const apiConfigName = "apiConfig";
async function getConfig(file) {
  const response = await fetch(file);
  const configTemp = await response.json();
  const IFrameMap = new Map();
  let i = 0;
  while (i < configTemp.Iframes.length) {
    // console.log(configTemp.Iframes[i]);
    IFrameMap.set(configTemp.Iframes[i].apiName, configTemp.Iframes[i]);
    i++;
  }
  const mapArray = Array.from(IFrameMap);
  const jsonObject = Object.fromEntries(mapArray);
  const jsonConfig = JSON.stringify(jsonObject);
  localStorage.setItem(apiConfigName, jsonConfig);
}
