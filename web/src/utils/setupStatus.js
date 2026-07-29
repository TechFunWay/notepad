import { getSetupStatus } from '@/api/auth'

let setupRequired = null
let setupCheckPromise = null

export async function checkSetupRequired() {
  if (setupRequired !== null) return setupRequired

  if (!setupCheckPromise) {
    setupCheckPromise = getSetupStatus()
      .then(({ data }) => {
        setupRequired = Boolean(data.needs_setup)
        return setupRequired
      })
      .finally(() => {
        setupCheckPromise = null
      })
  }

  return setupCheckPromise
}

export function markSetupComplete() {
  setupRequired = false
}
