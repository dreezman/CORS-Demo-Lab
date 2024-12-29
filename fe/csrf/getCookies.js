async function getCookies() {
  let cookiecount = 1;

  try {
    // try to clear all cookies
    document.cookie.split(";").forEach(function (c) {
      document.cookie =
        c.trim().split("=")[0] +
        "=;" +
        "expires=Thu, 01 Jan 1970 00:00:00 UTC;";
    });
    // fetch cookies
    const response = await fetch("/get-cookies", {
      method: "GET",
    });
    const data = await response.json();
    console.log(data);
  } catch (e) {
    if (!(e instanceof Error)) {
      e = new Error(e);
    }
    console.error(e.message);
  }
}

function openCookieJar(cookieObject) {
  document.cookie.split(";").forEach(function (c) {
    document.getElementById(cookieObject).innerHTML += c + "\n";
  });
}
