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

		var tagsContainer = document.getElementsByClassName('form-tag-list');

		if (tagsContainer.length === 0)
			return ;

		tagsContainer = tagsContainer[0];

		var tags = tagsContainer.getElementsByClassName('form-tag-item')
		for (var i = 0; i < tags.length; i++) {
			tags[i].getElementsByTagName('button')[0].onclick = removeTag;
		}

		var i = 0;
		var tagsInput = document.getElementsByClassName('form-tag-add')[0];

		tagsInput.onkeypress = function(event) {
			if (event.which === 13 || event.keyCode === 13) {
				event.preventDefault();
				return addTags(tagsInput.value);
			}
		};

		function addTags(value) {
			var value = tagsInput.value;
			if (value === undefined || value.length === 0)
				return ;
			tagsContainer.insertBefore(new buildElem(value), tagsInput);
			tagsInput.value = '';
		}

		function buildElem(value) {
			var span = document.createElement('span');
			var input = document.createElement('input');
			var removeButton = document.createElement('button');

			span.classList.add('form-tag-item');

			input.setAttribute('value', value);
			input.setAttribute('name', 'tags[]');
			input.setAttribute('type', 'text');
			input.setAttribute('readonly', '');

			removeButton.innerHTML = 'x';
			removeButton.onclick = removeTag;

			span.appendChild(input);
			span.appendChild(removeButton);
			return span;
		}

		function removeTag(event) {
			event.preventDefault();
			var p = this.parentNode;
			p.parentNode.removeChild(p);
		}

	});

})();
