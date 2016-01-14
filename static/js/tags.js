(function() {
	'use strict';

	// Copy|Delete|manage tags inputs
	document.addEventListener('DOMContentLoaded', function() {

		/*
		** Template
		*
		*	<input class="tags-add-input" /> -- Main input
		*	<div class="tags-container"> -- Tags container
		** 	<button class="tags-add-button"> -- If you also want an 'add' button (Cf Khaled)
		*/

		var tagsContainer = document.getElementsByClassName('tags-container');

		if (tagsContainer.length === 0)
			return ;

		tagsContainer = tagsContainer[0];

		var i = 0;
		var tagsInput = document.getElementsByClassName('tags-add-input')[0];

		tagsInput.onkeypress = function(event) {
			if (event.which === 13 || event.keyCode === 13) {
				event.preventDefault();
				return addTags(tagsInput.value);
			}
		};

		document.getElementsByClassName('tags-add-button')[0].onclick = function() {
			event.preventDefault();
			return addTags(tagsInput.value);
		}

		function addTags(value) {
			var value = tagsInput.value;
			if (value === undefined || value.length === 0)
				return ;
			tagsContainer.appendChild(new buildElem(value));
			tagsInput.value = '';
		}

		function buildElem(value) {
			var div = document.createElement('div');
			var input = document.createElement('input');
			var removeButton = document.createElement('button');

			input.value = value;
			input.name = 'tags[]';
			input.type = 'text';
			input.readonly = true;

			removeButton.innerHTML = 'Remove';
			removeButton.onclick = removeTag;

			div.appendChild(input);
			div.appendChild(removeButton);
			return div;
		}

		function removeTag() {
			event.preventDefault();
			var p = this.parentNode;
			p.parentNode.removeChild(p);
		}

	});

})();