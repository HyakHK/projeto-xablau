//default chatroom
var selectedchat= "general";


//update chatroom id
function changeChatRoom() {
    var newchat= document.getElementById("chatroom");
    if (newchat != null && newchat.value != selectedchat) {
        console.log(newchat);
    }
    return false;
}

//Send Message
function sendMessage() {
    var newmessage= document.getElementById("message");
    if (newmessage != null) {
        console.log(newmessage);
        newmessage.value = ""; //clear messag field after sending
    }
    return false;
}

//on load run

window.onload= function(){
    //avoid unecessary redirect
    document.getElementById("chatroom-selection").onsubmit= changeChatRoom;
    document.getElementById("chatroom-message").onsubmit= sendMessage;

    //check compatibility with WebSocket
    if (window["WebSocket"]) {
        console.log("supports websockets");
    } else {
        alert("Websocket not supported");
    }
};