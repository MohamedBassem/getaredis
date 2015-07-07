function fixLoaderPosition(){
  var loader = document.getElementById("loader");
  var loaderContainerWidth = loader.clientWidth;
  var loaderWidth = parseInt(loader.children[0].clientWidth)*10;
  loader.children[0].style.marginLeft = (loaderContainerWidth/2 - loaderWidth/2 - 10) + "px";
}

(function() {
  fixLoaderPosition()
  window.addEventListener("resize", fixLoaderPosition)
  
  var button = document.getElementById("getaredis");
  var errorMessage = document.getElementById("error");
  var message = document.getElementById("message");
  button.addEventListener("click", function() {
    button.style.display = "none";
    errorMessage.style.display = "none";
    message.style.display = "none";
    loader.style.visibility = "visible";
    var http = new XMLHttpRequest();
    var url = "/instance";
    http.open("POST", url, true);

    http.onreadystatechange = function() { //Call a function when the state changes.
      if(http.readyState == 4) {
        loader.style.visibility = "hidden";
        if(http.status == 200){
          console.log(http.responseText);
          data = JSON.parse(http.responseText);
          message.innerHTML = "IP: " + data["IP"] + "<br/>Port: " + data["port"] + "<br/>Password: " + data["password"] + "<br/>";
          message.innerHTML+= "Example:<br/>$ telnet " + data["IP"] + " " + data["port"] + "<br/>AUTH " + data["password"] +"<br/>";
          message.innerHTML+= "+OK<br/>PING<br/>+PONG<br/>";
          message.style.display = "";
        }else{
          errorMessage.innerHTML = http.responseText;
          errorMessage.style.display = "";
          button.style.display = "";
        }
      }
    }
    http.send();
  })

})();
