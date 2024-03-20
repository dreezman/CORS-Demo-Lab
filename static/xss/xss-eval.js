function getUrlParameter(name) {
    return (location.search.split(name + '=')[1] || '').split('&')[0]
}
document.cookie="TopSecretPassword=1234567890";
var op1 = decodeURI(getUrlParameter('op1'))
var op2 = decodeURI(getUrlParameter('op2'))
var sum = eval(`${op1} + ${op2}`)
console.log(`The sum is: ${sum}`) //Inject XSS into the console
