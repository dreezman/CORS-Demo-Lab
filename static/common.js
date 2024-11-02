
// when document is loaded, find the parameters on the HTTP call
/*
Used like this
var op1 = decodeURI(getUrlParameter('op1'))
var op2 = decodeURI(getUrlParameter('op2'))
*/
function getUrlParameter(name) {
    return (location.search.split(name + '=')[1] || '').split('&')[0]
}


getConfig("../config.json"); // data is stored in localstorage under "config"
const jsonConfig = JSON.parse(localStorage.getItem("iframeConfig"));

