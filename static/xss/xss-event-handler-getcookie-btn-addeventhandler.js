function getCookie(){
    document.cookie="TopSecretPassword=1234567890";
    alert('I see your cookies!!!:   ' + document.cookie)
}
document.addEventListener("DOMContentLoaded", function() {
    document.getElementById("myBtn").addEventListener("click", getCookie);
});

