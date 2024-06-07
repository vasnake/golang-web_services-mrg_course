function encodeHTML(s) {
    return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/"/g, '&quot;');
}

function NewRequest(method, url) {
    var request = new XMLHttpRequest();
    request.open(method, url, true);
    var csrf_token = document.querySelector("meta[name='csrf-token']").getAttribute("content");
    request.setRequestHeader("csrf-token", csrf_token);

    return request;
}

function NewGQLRequest() {
    return NewRequest("POST", "/graphql")
}

const rateCommentToggleMuration = `
mutation($photoID: ID!, $direction: String!) {
    ratePhoto(photoID:$photoID, direction: $direction){
        id
        rating
    }
}
`

function rateCommentToggle(elem) {
    var voteDir = "up";
    if(elem.classList.contains('hi-red')) {
        elem.classList.remove('hi-red');
        voteDir = "down";
    } else {
        elem.classList.add('hi-red');
    }

    var request = NewGQLRequest();
    request.setRequestHeader('Content-Type', 'application/json');
    var params = {
        variables: {
            photoID: elem.getAttribute('data-id'),
            direction: voteDir,
        },
        query: rateCommentToggleMuration,
    };
    var body = JSON.stringify(params);
    
    request.onload = function() {
        var resp = JSON.parse(request.responseText);
        if(resp.errors) {
            console.log("rateCommentToggle server err:", resp.errors);
            return;
        }
        var ratingElem = document.querySelector('#rating-'+resp.data.ratePhoto.id);
        // rating = parseInt(ratingElem.innerHTML) + parseInt(vote);
        ratingElem.innerHTML = resp.data.ratePhoto.rating;
    };
    request.send(body);
}

  
const getPhotosQuery = `query($userID: ID!) {
    user(userID: $userID) {
      id
      name
      avatar
      photos {
        id
        user {id, name, avatar, followed}
        url
        comment
        rating
        liked
      }
    }
    me {
        id
        name
        avatar
        followedUsers {id, name, avatar, followed}
        recomendedUsers {id, name, avatar, followed}
    }
  }
`

function getUserPhotos(uid) {
    var request = NewGQLRequest();

    request.setRequestHeader('Content-Type', 'application/json');
    var params = {
        variables: {
            userID: uid
        },
        query: getPhotosQuery,
    };
    var body = JSON.stringify(params);
    
    request.onload = function() {
        var resp = JSON.parse(request.responseText);
        if(resp.errors) {
            console.log("getUserPhotos server err:", resp.errors);
            return;
        }

        renderPhotosHTML(resp.data.user.photos);
        renderUserListHTML(resp.data.me.followedUsers, "following");
        renderUserListHTML(resp.data.me.recomendedUsers, "recomends");
    };
    request.onerror = function() {
        console.log("renderPhotos error", request.responseText)
    }
    request.send(body);
}

function renderPhotosHTML(elems) {
    msgs = document.getElementById("photolist");
    msgs.innerHTML = "";
    elems.forEach(elem => {
        // you better to use modern JS frameworks in production
        msgNode = `<div>
            <div style="border-bottom:1px solid silver; padding:4px; font-size:14px;">
                <a href="/photos/${elem.user.name}" class="userName"><img style="border-radius: 15px;" height=32 width=32 src="${elem.user.avatar}" /> ${elem.user.name}</a>
                <a onclick="followUser(this);" data-id="${elem.user.id}" href="#">${elem.user.followed ? "[unfollow]" : "[follow]"}</a>
            </div>
            <img src="/images/${elem.user.id}/${elem.url}_600.jpg" />
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
}

/*
curl localhost:8080/query \
  -H 'Cookie: session=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjYsImV4cCI6MTU3OTM1MzQ0OCwianRpIjoieVNOempIRGFYWWFHQVJCS2ljaENiQWFpS2RFT3JuY2MiLCJpYXQiOjE1NzE1Nzc0NDh9.iDp_yr9Qhd5LXnOM1Ocvhkhp6u27X7jLtPTmrFGZOqk' \
  -F operations='{ "query": "mutation($comment: String!, $file: Upload!) { uploadPhoto(comment: $comment, file: $file) { id } }", "variables": { "comment": "uploaded by graphql", "file": null } }' \
  -F map='{ "0": ["variables.file"] }' \
  -F 0=@./photo_samples/building_1.jpg \
  --trace-ascii -

это сгенерирует вот такой запрос

------WebKitFormBoundaryyxMrNgnvhVtaooUv
Content-Disposition: form-data; name="my_file"; filename="building_9.jpg"
Content-Type: image/jpeg


------WebKitFormBoundaryyxMrNgnvhVtaooUv
Content-Disposition: form-data; name="operations"

{"query":"\nmutation($comment: String!, $file: Upload!) { \n    uploadPhoto(comment: $comment, file: $file) { \n        id, comment\n    }\n}\n","variables":{"comment":"building 9 comment","file":null}}
------WebKitFormBoundaryyxMrNgnvhVtaooUv
Content-Disposition: form-data; name="map"

{"my_file":["variables.file"]}
------WebKitFormBoundaryyxMrNgnvhVtaooUv--
*/

const uploadPhotoMutation = `
mutation($comment: String!, $file: Upload!) { 
    uploadPhoto(comment: $comment, file: $file) { 
        id, comment
    }
}
`

function uploadPhoto(uid) {
    var form = new FormData(document.getElementById('uploadPhoto'));
    var form2send = new FormData();

    form2send.set("operations", JSON.stringify({
        query: uploadPhotoMutation,
        variables: {
            comment: form.get("comment"),
            file: null,
        }
    }));
    form2send.set("map", JSON.stringify({
        my_file: ["variables.file"]
    }));
    form2send.set("my_file", form.get("my_file"));

    var request = NewGQLRequest();
    request.onload = function() {
        var resp = JSON.parse(request.responseText);
        if(resp.errors) {
            console.log("uploadPhoto server err:", resp.errors);
            return;
        }
        getUserPhotos(uid);
    };
    request.send(form2send);
}

function renderUserListHTML(data, selector) {
    msgs = document.getElementById(selector);
    msgs.innerHTML = "";
    data.forEach(elem => {
        // you better to use modern JS frameworks in production
        msgNode = `
            <a href="/photos/${elem.name}" class="userName"> ${elem.name}</a>
            <a href="#" onclick="followUser(this);" data-id="${elem.id}">${elem.followed ? "[unfollow]" : "[follow]"}</a>
        `;

        var node = document.createElement("div");
        node.className = "userElem";
        node.innerHTML = msgNode;

        msgs.appendChild(node);
    });
}

const myFollowersQuery = `
query {
    me {
      id
      name
      avatar
      followedUsers {id, name, avatar, followed}
      recomendedUsers {id, name, avatar, followed}
    }
}
`

function getUserList() {
    var request = NewGQLRequest();

    request.setRequestHeader('Content-Type', 'application/json');
    var params = {
        query: myFollowersQuery,
    };
    var body = JSON.stringify(params);

    request.onload = function() {
        var resp = JSON.parse(request.responseText);
        if(resp.errors) {
            console.log("getUserList server err:", resp.errors);
            return;
        }
        renderUserListHTML(resp.data.me.followedUsers, "following");
        renderUserListHTML(resp.data.me.recomendedUsers, "recomends");
    };
    request.onerror = function() {
        console.log("renderUserList error", request.responseText)
    }
    request.send(body);
}

const followUserMutation = `
mutation($userID: ID!, $direction: String!) {
    followUser(userID:$userID, direction:$direction){
        id
        name
        avatar
    }
}
`

function followUser(elem) {
    var request = NewGQLRequest();
    request.setRequestHeader('Content-Type', 'application/json');
    var params = {
        variables: {
            userID: elem.getAttribute('data-id'),
            direction: elem.innerHTML == "[unfollow]" ? "down" : "up",
        },
        query: followUserMutation,
    };
    var body = JSON.stringify(params);

    request.onload = function() {
        var resp = JSON.parse(request.responseText);
        if(resp.errors) {
            console.log("followUser server err:", resp.errors);
            return;
        }
        getUserList()
    };
    request.send(body);
    return false;
}