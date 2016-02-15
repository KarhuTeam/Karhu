(function() {
	'use strict';

	document.addEventListener('DOMContentLoaded', function() {

		/*
		** Template
		*
        *   <ul class="fold-container">
        *       <li class="fold-parent" data-fold-id="1"></li>
        *       <li class="fold-child fold-child-hidden" data-fold-id="1"></li>
        *       <li class="fold-parent" data-fold-id="2"></li>
        *       <li class="fold-child fold-child-hidden" data-fold-id="2"></li>
        *   </ul>
		*/

		var foldContainer = document.getElementsByClassName('fold-container');

		if (foldContainer.length === 0)
			return ;

		for (var i = 0; i < foldContainer.length; i++) {

			new FoldContainer(foldContainer[i]);
		}

		function FoldContainer(element) {

			this.foldContainer = element;
            var self = this; // #Hacker

			var folds = this.foldContainer.getElementsByClassName('fold-parent');
			for (var i = 0; i < folds.length; i++) {

                folds[i].onclick = function(event) {

					this.classList.toggle('fold-parent-selected');

                    var id = this.getAttribute('data-fold-id');
                    var childs = self.foldContainer.getElementsByClassName('fold-child');

                    for ( var j = 0; j < childs.length; j++) {
                        if (childs[j].getAttribute('data-fold-id') == id) {
                            childs[j].classList.toggle('fold-child-hidden');
                        }
                    }
                };
			}
		}
	});

})();
