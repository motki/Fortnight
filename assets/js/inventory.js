;(function ($, window, undefined) {
  const $loading = $(document.getElementById('loading'));
  const $overlay = $(document.getElementById('overlay'));
  const $error = $(document.getElementById('error'));
  const $errorReason = $(document.getElementById('error-reason-text'));
  const $mainWindow = $(document.getElementById('main-window'));
  const $locations = $(document.getElementById('locations'));
  const $items = $(document.getElementById('items'));

  // Show an extra message if loading takes forever.
  let ti = setTimeout(() => {
    ti = undefined;
    if ($loading.is(':visible')) {
      $loading.find('#loading-message').fadeOut('fast', function() {
        $(this).text('Still working on it...').fadeIn('fast');
      });
    }
  }, 5000);

  $mainWindow.on('motki:loaded', () => {
    if (ti !== undefined) {
      clearTimeout(ti);
      ti = undefined;
    }
    $loading.fadeOut('fast');
    $overlay.fadeOut('fast');
  });

  $mainWindow.on('motki:error', (e, err) => {
    let message = err.message || 'unknown';
    if (err.context) {
      message += ` ${err.context}`;
    }
    console.debug(message);
    $errorReason.text(message);
    if (!$overlay.is(':visible')) {
      $overlay.fadeIn('fast');
    }
    if ($loading.is(':visible')) {
      $loading.fadeOut('fast', function() {
        $error.fadeIn('fast');
      });
    } else {
      $error.fadeIn('fast');
    }
  });

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

    itemBelowThreshold: '<i class="icon icon-alert item-low-alert"></i>'
  };

  class InventoryItem {
    constructor(props) {
      this.typeId = props['type_id'];
      this.locationId = props['location_id'];
      this.currentLevel = props['current_level'];
      this.minimumLevel = props['minimum_level'];
      this.fetchedAt = moment(props['fetched_at']);
      this.name = props['name'];
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

      $locations.on('click', '.nav-group-item', (e) => {
        e.preventDefault();
        console.log(e);
        this.selectLocation($(e.currentTarget).data('id'));
      });
    }

    fetch() {
      return fetch('/inventory').then((its) => {
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
      });
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
  }

  $(function () {
    const inventory = new Inventory();

    inventory.fetch().catch(handleError);
  });
})($, window);