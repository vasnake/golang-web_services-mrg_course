function rateCommentToggle(elem) {
    var id = elem.getAttribute('data-id');
    var vote = 0;
    if(elem.classList.contains('hi-red')) {
        // remove user like
        elem.classList.remove('hi-red');
        vote = -1;
    } else {
        // add user like
        elem.classList.add('hi-red');
        vote = 1;
    }

    var csrf_token = document.querySelector("meta[name='csrf-token']").getAttribute("content");

    var request = new XMLHttpRequest();
    request.open('POST', '/photos/rate?id='+id+"&vote="+(vote > 0 ? "up" : "down"), true);
    request.setRequestHeader("csrf-token", csrf_token);
    request.onload = function() {
        var resp = JSON.parse(request.responseText);
        if(resp.err) {
            console.log("rateComment server err:", resp.err);
            return;
        }
        var ratingElem = document.querySelector('#rating-'+resp.id);
        rating = parseInt(ratingElem.innerHTML) + parseInt(vote); // rating not live, it's snapshot in time of open list page
        ratingElem.innerHTML = rating;
    };

    request.send();
}
