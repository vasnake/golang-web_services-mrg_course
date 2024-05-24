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
    request.open('POST', '/photos/rate?id='+id+"&vote="+(vote > 0 ? "up" : "down"), true);
    
    var csrf_token = document.querySelector("meta[name='csrf-token']").getAttribute("content");
    
    request.setRequestHeader("csrf-token", csrf_token);
    request.onload = function() {
        var resp = JSON.parse(request.responseText);
        if(resp.err) {
            console.log("rateComment server err:", resp.err);
            return;
        }
        var ratingElem = document.querySelector('#rating-'+resp.id);
        rating = parseInt(ratingElem.innerHTML) + parseInt(vote);
        ratingElem.innerHTML = rating;
    };
    request.send();
}

