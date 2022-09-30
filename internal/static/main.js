
const CONTAINER_STATE_STALE = 'Stale'
const CONTAINER_STATE_FRESH = 'Fresh'
const CONTAINER_STATE_REQUESTED_UPDATE = 'RequestedUpdate'

document.addEventListener('DOMContentLoaded', () => {
  const buttons = document.querySelectorAll('.container-info > a.button')
  buttons.forEach(currentBtn => {
    currentBtn.addEventListener('click', event => {
      event.preventDefault()
      updateButtonClick(event)
    }, false)
  })

  const updateAllButton = document.querySelector('#update-all')
  if (updateAllButton !== null) {
    updateAllButton.addEventListener('click', event => {
      event.preventDefault()
      updateAllButtonClick(event)
    }, false)
  }

  const forceRefreshButton = document.querySelector('#force-refresh')
  if (forceRefreshButton !== null) {
    forceRefreshButton.addEventListener('click', event => {
      event.preventDefault()
      forceRefreshButtonClick(event)
    }, false)
  }

  mainPoll()
})

function mainPoll() {
  pollContainers().then(() => {
    setTimeout(mainPoll, 1000)
  })
}

async function pollContainers() {
  const query = new URLSearchParams({
    after: window.LAST_UPDATE,
  })
  const res = await fetch(`/api/v1/containers?${query}`)
  if (res.status === 304) {
    return
  }

  const info = await res.json()

  window.LAST_UPDATE = info.LastUpdate
  const containers = info.Containers

  const infos = Array.from(document.querySelectorAll('.container-info'))

  if (infos.length !== containers.length) {
    location.reload()
  }

  for (let info of infos) {
    const current = containers.find(el => el.ID === info.dataset.id)
    if (current === undefined) {
      // no container info in response
      info.dataset.state = CONTAINER_STATE_FRESH
      continue
    }

    if (info.dataset.state !== current.State)
      info.dataset.state = current.State
  }
}

async function updateButtonClick(event) {
  const parent = event.target.closest('.container-info')
  const containerId = parent.dataset.id

  if (parent.dataset.state === CONTAINER_STATE_STALE) {
    parent.dataset.state = CONTAINER_STATE_REQUESTED_UPDATE

    await fetch('/api/v1/containers/update', {
      method: 'POST',
      body: JSON.stringify({ Ids: [containerId] })
    })
  }
}

async function updateAllButtonClick(event) {
  let ids = []
  const infos = Array.from(document.querySelectorAll(`.container-info[data-state="${CONTAINER_STATE_STALE}"]`))
  for (let info of infos) {
    ids.push(info.dataset.id)
    info.dataset.state = CONTAINER_STATE_REQUESTED_UPDATE
  }

  if (ids.length > 0) {
    await fetch('/api/v1/containers/update', {
      method: 'POST',
      body: JSON.stringify({ Ids: ids })
    })
  }
}

async function forceRefreshButtonClick(event) {
  await fetch('/v1/update')
  event.target.style.pointerEvents = 'none'
}
