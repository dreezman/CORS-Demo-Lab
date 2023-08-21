async function getConfig(file) {
  const response = await fetch(file);
  const config = await response.json();
  const jsonConfig = JSON.stringify(config);
  localStorage.setItem("config", jsonConfig);
}
