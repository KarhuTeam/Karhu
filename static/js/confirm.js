(function() {
  'use strict';

  document.addEventListener('DOMContentLoaded', function () {
    var elements = document.getElementsByClassName('btn-danger');
    if (elements !== null) {
        for (var i = 0; i < elements.length; i++) {
            elements[i].addEventListener("click", function (e) {
                var valid = confirm(this.getAttribute("data-text"));
                if (!valid) {
                    e.preventDefault();
                }
            });
        }
    }
  });
})();
