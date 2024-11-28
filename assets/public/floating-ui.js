import {computePosition, autoUpdate, flip, offset} from 'https://cdn.jsdelivr.net/npm/@floating-ui/dom@1.6.12/+esm';

let updateMap = {};

function updatePosition(ref, floating) {
    let opts = {
        middleware: [flip(), offset(5)],
        placement: 'bottom-start',
    }

    updateMap[ref] = autoUpdate(ref, floating, () => {
        computePosition(ref, floating, opts).then((update) => {
            floating.style.left = update.x + 'px';
            floating.style.top = update.y + 'px';
        });
    })
}

function stopUpdate(ref) {
    if(updateMap[ref]) {
        updateMap[ref]()
        delete updateMap[ref];
    }
}

window.dockman = window.dockman || {}

window.dockman.floating = {
    updateMap,
    updatePosition,
    stopUpdate
}
