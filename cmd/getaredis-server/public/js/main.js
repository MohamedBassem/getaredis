function fixLoaderPosition(){
  var loader = document.getElementById("loader");
  var loaderContainerWidth = loader.clientWidth;
  var loaderWidth = parseInt(loader.children[0].clientWidth)*10;
  loader.children[0].style.marginLeft = (loaderContainerWidth/2 - loaderWidth/2 - 10) + "px";
}

function fixFooterPosition(){
  var height = window.innerHeight;
  var container = document.getElementById("container").parentNode;
  var footerHeight = document.getElementById("footer").clientHeight;
  if( container.clientHeight + footerHeight > height ){
    container.style.marginBottom = "0px";
  }else{
    container.style.marginBottom = (height-container.clientHeight-footerHeight) + "px";
  }
}

(function() {
  fixLoaderPosition();
  fixFooterPosition();
  window.addEventListener("resize", fixLoaderPosition);
  window.addEventListener("resize", fixFooterPosition);
  
  var button = document.getElementById("getaredis");
  var errorMessage = document.getElementById("error");
  var message = document.getElementById("message");
  button.addEventListener("click", function() {
    button.style.display = "none";
    errorMessage.style.display = "none";
    message.style.display = "none";
    loader.style.visibility = "visible";
    fixFooterPosition();
    var http = new XMLHttpRequest();
    var url = "/instance";
    http.open("POST", url, true);

    http.onreadystatechange = function() { //Call a function when the state changes.
      if(http.readyState == 4) {
        loader.style.visibility = "hidden";
        if(http.status == 200){
          data = JSON.parse(http.responseText);
          message.innerHTML = '<strong style="color:grey;"># Your Instance is Ready!</strong>\n'
          message.innerHTML+= "$ telnet " + data["IP"] + " " + data["port"] + "\n";
          message.innerHTML+= "AUTH " + data["password"] +"\n";
          message.innerHTML+= '<span style="color:yellow;">+OK</span>\n';
          message.style.display = "";
        }else{
          errorMessage.innerHTML = http.responseText;
          errorMessage.style.display = "";
          button.style.display = "";
        }
        fixFooterPosition();
      }
    }
    http.send();
  })

})();
