function encodeHTML(s) {
    return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/"/g, '&quot;');
}

function rateCommentToggle(elem) {
    var id = elem.getAttribute('data-id');
    var vote = 0;
    if(elem.classList.contains('hi-red')) {
        elem.classList.remove('hi-red');
        vote = -1;
    } else {
        elem.classList.add('hi-red');
        vote = 1;
    }

    var request = new XMLHttpRequest();
    request.open('POST', '/api/v1/photos/rate?id='+id+"&vote="+(vote > 0 ? "up" : "down"), true);
    
    var csrf_token = document.querySelector("meta[name='csrf-token']").getAttribute("content");
    
    request.setRequestHeader("csrf-token", csrf_token);
    request.onload = function() {
        var resp = JSON.parse(request.responseText);
        if(resp.error) {
            console.log("rateComment server err:", resp.err);
            return;
        }
        var ratingElem = document.querySelector('#rating-'+resp.body.id);
        rating = parseInt(ratingElem.innerHTML) + parseInt(vote);
        ratingElem.innerHTML = rating;
    };
    request.send();
}

function renderPhotos(uid) {
    var request = new XMLHttpRequest();
    request.open('GET', '/api/v1/photos/list?uid='+uid, true);
    
    var csrf_token = document.querySelector("meta[name='csrf-token']").getAttribute("content");
    
    request.setRequestHeader("csrf-token", csrf_token);
    request.onload = function() {
        var resp = JSON.parse(request.responseText);
        if(resp.error) {
            console.log("renderPhotos server err:", resp.err);
            return;
        }

        msgs = document.getElementById("photolist");
        msgs.innerHTML = "";
        resp.body.photolist.forEach(elem => {
            // you better to use modern JS frameworks in production
            msgNode = `<div>
                <div style="border-bottom:1px solid silver; padding:4px; font-size:14px;">
                    <a href="/photos/${elem.user_login}" class="userName"><img style="border-radius: 15px;" src="/images/${elem.path}_32.jpg" /> ${elem.user_login}</a>
                    <a onclick="followUser(this);" data-id="${elem.user_id}" href="#">${elem.followed ? "[unfollow]" : "[follow]"}</a>
                </div>
				<img src="/images/${elem.path}_600.jpg" />
				<div class="details">
					<span onclick="rateCommentToggle(this)" data-id="${elem.id}" class="hi ${elem.liked === true ? 'hi-red' : ''}">❤</span>
					<span class="rating" id="rating-${elem.id}">${elem.rating}</span>
					<br/>
					${encodeHTML(elem.comment)}
                </div>
                <div class="commentForm">
                    <form onsubmit="return false">
                        <input type="text" placeholder="Comment text...">
                        <input type="submit" value="Comment">
                    </form>
                </div>
			</div>`;
    
            var node = document.createElement("div");
            node.className = "photoElem";
            node.innerHTML = msgNode;
    
            msgs.appendChild(node);
        });
    };
    request.onerror = function() {
        console.log("renderPhotos error", request.responseText)
    }
    request.send();
}

function uploadPhoto(uid) {
    var form = new FormData(document.getElementById('uploadPhoto'))
    var request = new XMLHttpRequest();
    request.open('POST', '/api/v1/photos/upload', true);
    var csrf_token = document.querySelector("meta[name='csrf-token']").getAttribute("content");
    request.setRequestHeader("csrf-token", csrf_token);
    request.onload = function() {
        var resp = JSON.parse(request.responseText);
        if(resp.error) {
            console.log("rateComment server err:", resp.err);
            return;
        }
        renderPhotos(uid);
    };
    request.send(form);
}

function renderUserList(list) {
    var request = new XMLHttpRequest();
    request.open('GET', '/api/v1/user/' + list, true);
    var csrf_token = document.querySelector("meta[name='csrf-token']").getAttribute("content");
    request.setRequestHeader("csrf-token", csrf_token);
    request.onload = function() {
        var resp = JSON.parse(request.responseText);
        if(resp.error) {
            console.log("renderUserList server err:", resp.err);
            return;
        }

        msgs = document.getElementById(list);
        msgs.innerHTML = "";
        resp.body.users.forEach(elem => {
            // you better to use modern JS frameworks in production
            msgNode = `
                <a href="/photos/${elem.login}" class="userName"> ${elem.login}</a>
                <a href="#" onclick="followUser(this);" data-id="${elem.id}">${resp.body.followed ? "[unfollow]" : "[follow]"}</a>
			`;
    
            var node = document.createElement("div");
            node.className = "userElem";
            node.innerHTML = msgNode;
    
            msgs.appendChild(node);
        });
    };
    request.onerror = function() {
        console.log("renderUserList error", request.responseText)
    }
    request.send();
}

function followUser(elem) {
    var id = elem.getAttribute('data-id');
    var unfollow = elem.innerHTML == "[unfollow]" ? 1 : 0;
    var request = new XMLHttpRequest();
    request.open('POST', '/api/v1/user/follow?id='+id+"&unfollow="+unfollow, true);
    var csrf_token = document.querySelector("meta[name='csrf-token']").getAttribute("content");
    request.setRequestHeader("csrf-token", csrf_token);
    request.onload = function() {
        var resp = JSON.parse(request.responseText);
        if(resp.error) {
            console.log("rateComment server err:", resp.err);
            return;
        }
        renderUserList("following");
        renderUserList("recomends");
    };
    request.send();
    return false;
}