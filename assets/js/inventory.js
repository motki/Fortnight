;(function ($, window, undefined) {
  const $loading = $(document.getElementById('loading'));
  const $overlay = $(document.getElementById('overlay'));
  const $error = $(document.getElementById('error'));
  const $errorReason = $(document.getElementById('error-reason-text'));
  const $mainWindow = $(document.getElementById('main-window'));
  const $locations = $(document.getElementById('locations'));
  const $items = $(document.getElementById('items'));
  const $refresh = $(document.getElementById('refresh-data'));
  const $newItem = $(document.getElementById('new-item'));
  const $edit = $(document.getElementById('edit'));
  const $editForm = $(document.getElementById('edit-form'));
  const $editCancel = $(document.getElementById('cancel-edit'));
  const $editSave = $(document.getElementById('save-edit'));

  let currentModal;

  const openModal = ($el) => {
    return new Promise((resolve) => {
      if (currentModal !== undefined) {
        currentModal.fadeOut('fast');
      }
      currentModal = $el;
      $overlay.fadeIn('fast');
      $el.fadeIn('fast');
      return resolve();
    });
  };

  const hideModal = () => {
    return new Promise((resolve) => {
      if (currentModal !== undefined) {
        currentModal.fadeOut('fast');
        $overlay.fadeOut('fast');
        currentModal = undefined;
      }
      return resolve();
    })
  };

  // Show an extra message if loading takes forever.
  let tiFactory = () => setTimeout(() => {
    ti = undefined;
    if ($loading.is(':visible')) {
      $loading.find('#loading-message').fadeOut('fast', function() {
        $(this).text('Still working on it...').fadeIn('fast');
      });
    }
  }, 5000);
  
  let ti;

  $mainWindow.on('motki:loading', () => {
    $loading.find('#loading-message').text('Please wait...');
    openModal($loading);
    ti = tiFactory();
  });

  $mainWindow.on('motki:loaded', () => {
    if (ti !== undefined) {
      clearTimeout(ti);
      ti = undefined;
    }
    hideModal();
  });

  $mainWindow.on('motki:error', (e, err) => {
    let message = err.message || 'unknown';
    if (err.context) {
      message += ` ${err.context}`;
    }
    console.debug(message);
    $errorReason.text(message);
    openModal($error);
  });

  $mainWindow.on('motki:new-item', () => {
    if ($overlay.is(':visible')) {
      return;
    }
    $editForm.html(markup.itemForm(new InventoryItem()));
    openModal($edit);
  });

  $mainWindow.on('motki:confirm', () => {
    if (!$editForm.is(':visible')) {
      return;
    }
    $editForm.find('form').submit();
    hideModal();
  });

  $mainWindow.on('motki:cancel', () => {
    if (!$editForm.is(':visible')) {
      return;
    }
    hideModal();
  });

  $editCancel.on('click', () => $mainWindow.trigger('motki:cancel'));
  $editSave.on('click', () => $mainWindow.trigger('motki:confirm'));

  const handleError = (err) => {
    let e, ctx;
    if (err instanceof Array) {
      [e, ctx] = err;
    } else {
      e = err;
    }
    console.debug('handleError', e, ctx);
    $mainWindow.trigger('motki:error', {message: e, context: ctx});
  };

  window.addEventListener('error', (e) => {
    console.debug('window.onerror', e);
    handleError(e.error || e.message);
    return false;
  });

  const fetch = (url, context) => {
    const ctx = context || `fetching ${url}`;
    return new Promise(function (resolve, reject) {
      $.ajax({
        url: url,
        dataType: 'json',
        success: (result) => {
          if (result.success) {
            resolve(result.data);
          } else {
            reject([result.data, ctx]);
          }
        },
        error: (resp) => {
          console.debug(resp);
          if (resp.responseJSON) {
            reject([resp.responseJSON.data, ctx]);
          } else {
            reject([resp.responseText || resp.statusText, ctx]);
          }
        },
        timeout: 30000
      });
    });
  };

  const groupBy = (items, prop) => {
    let res = {};
    for (let it of items) {
      const k = it[prop];
      if (res[k] === undefined) {
        res[k] = [it];
      } else {
        res[k].push(it);
      }
    }
    return res;
  };

  const markup = {
    /**
     * @param {Location} loc
     */
    locationNavItem:
      (loc) => `<span id="location-${loc.id}" data-id="${loc.id}" class="nav-group-item">` +
                 `${loc.name}` +
               `</span>`,

    /**
     * @param {InventoryItem} it
     */
    itemRow:
      (it) => `<tr>` +
                `<td>${it.typeId}</td>` +
                `<td>${it.name}</td>` +
                `<td class="item-low">${it.belowThreshold ? markup.itemBelowThreshold : ''}</td>` +
                `<td>${it.currentLevel}</td>` +
                `<td>${it.minimumLevel}</td>` +
                `<td>${it.fetchedAt.fromNow()}</td>` +
              `</tr>`,

    itemBelowThreshold: '<i class="icon icon-alert item-low-alert"></i>',

    itemForm:
      (it) => `<form method="post" action="/inventory">` +
                `<div class="form-group">` +
                  `<label>Location</label>` +
                  `<input type="text" class="form-control" id="loc-lookup">` +
                `</div>` +
                `<div class="form-group">` +
                  `<label>Item Type</label>` +
                  `<input type="text" class="form-control" id="typ-lookup">` +
                `</div>` +
                `<div class="form-group">` +
                  `<label>Low Level Threshold</label>` +
                  `<input type="text" class="form-control" name="minimumLevel" value="${it.minimumLevel || ''}">` +
                `</div>` +
                `<div class="form-group">` +
                  `<label>Current Level</label>` +
                  `<input type="text" class="form-control" name="currentLevel" disabled readonly value="${it.currentLevel || ''}">` +
                `</div>` +
                `<div class="form-group">` +
                  `<label>Last Updated</label>` +
                  `<input type="text" class="form-control" name="fetchedAt" disabled readonly value="${it.fetchedAt || ''}">` +
                `</div>`
  };

  class InventoryItem {
    constructor(props) {
      const p = props || {};
      this.typeId = p['type_id'];
      this.locationId = p['location_id'];
      this.currentLevel = p['current_level'];
      this.minimumLevel = p['minimum_level'];
      this.fetchedAt = moment(p['fetched_at']);
      this.name = p['name'];
      Object.defineProperty(this, 'belowThreshold', {
        get: () => this.minimumLevel > this.currentLevel
      });
    }
  }

  class Location {
    constructor(id, items) {
      this.id = id.toString();
      this.structure = undefined;
      this.station = undefined;
      this.system = undefined;
      this.items = undefined;
      if (items !== undefined) {
        this.setItems(items);
      }
      Object.defineProperty(this, 'name', {
        get: () => {
          if (this.system === undefined) {
            throw 'location data not populated';
          }
          if (this.structure !== undefined) {
            return this.structure.name;
          } else if (this.station !== undefined) {
            return this.station.name;
          } else if (this.system !== undefined) {
            return this.system.name;
          }
          return 'Unknown';
        }
      });
    }

    hydrate(propsOrItems) {
      if (propsOrItems instanceof Array) {
        this.setItems(propsOrItems);
      } else {
        this.structure = propsOrItems['structure'];
        this.station = propsOrItems['station'];
        this.system = propsOrItems['system'];
      }
    }

    setItems(items) {
      this.items = items.sort((a, b) => {
        if (a.belowThreshold) {
          if (b.belowThreshold) {
            return a.typeId < b.typeId ? -1 : 1;
          }
          return -1;
        }
        if (b.belowThreshold) {
          return 1;
        }
        return a.typeId < b.typeId ? -1 : 1;
      });
    }
  }

  class Inventory {
    constructor() {
      this.locationsById = {};
      this.inventoryByLocationId = {};
      this.currentLocation = undefined;
    }

    fetch() {
      return this.clear().then(fetch('/inventory').then((its) => {
        const mapped = its.map((it) => new InventoryItem(it));
        console.log(mapped);
        const inventoryByLocationId = groupBy(mapped, 'locationId');
        console.debug(this.inventoryByLocationId);
        const fetches = [];
        for (const [locID, items] of Object.entries(inventoryByLocationId)) {
          const location = this.locationsById[locID] = new Location(locID, items);
          fetches.push(fetch(`/location/${locID}`)
            .then(loc => {
              location.hydrate(loc);
              $locations.append(markup.locationNavItem(location));
              if (!this.currentLocation) {
                return this.selectLocation(location)
                  .then(() => $mainWindow.trigger('motki:loaded'));
              }
            }));
        }
        return Promise.all(fetches);
      })).catch(handleError);
    }

    getLocation(locationId) {
      return this.locationsById[locationId.toString()];
    }

    selectLocation(locationOrId) {
      return new Promise((resolve, reject) => {
        const lastLocId = this.currentLocation ? this.currentLocation.id : undefined;
        if (locationOrId instanceof Location) {
          this.currentLocation = locationOrId;
        } else {
          this.currentLocation = this.getLocation(locationOrId);
        }
        if (!this.currentLocation) {
          console.debug('current location is null');
          return reject(['location is invalid', 'selecting location']);
        } else if (lastLocId === this.currentLocation.id) {
          // Location is the same, return.
          return resolve();
        }
        $locations.find('.active').removeClass('active');
        $items.html('');
        $locations.find(`[data-id="${this.currentLocation.id}"]`).addClass('active');
        this.currentLocation.items.forEach(it => $items.append(markup.itemRow(it)));
        resolve();
      });
    }

    purge() {
      return fetch('/inventory/purge', `refreshing data`).then(() => this.fetch());
    }

    clear() {
      return new Promise((resolve) => {
        this.locationsById = {};
        this.inventoryByLocationId = {};
        this.currentLocation = undefined;
        $locations.html('');
        $items.html('');
        return resolve();
      });
    }
  }

  $(() => {
    const inventory = new Inventory();

    $locations.on('click', '.nav-group-item', (e) => {
      e.preventDefault();
      console.log(e);
      inventory.selectLocation($(e.currentTarget).data('id')).catch(handleError);
    });

    $refresh.on('click', () => {
      $mainWindow.trigger('motki:loading');
      inventory.purge().catch(handleError);
    });

    $newItem.on('click', () => {
      $mainWindow.trigger('motki:new-item');
    });

    $mainWindow.trigger('motki:loading');
    inventory.fetch().catch(handleError);
  });

  $(document).on('keyup', (e) => {
    switch (e.keyCode) {
      case 13:
        $mainWindow.trigger('motki:confirm');
        break;
      case 27:
        $mainWindow.trigger('motki:cancel');
        break;
      default:
        return true;
    }
    e.preventDefault();
    return false;
  });
})($, window);