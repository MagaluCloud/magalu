run "validate_active_engines" {
  command = apply

  assert {
    condition = alltrue([
      for engine in data.mgc_dbaas_engines.active_engines.engines : engine.status == "ACTIVE"
    ])
    error_message = "Found engines that are not in ACTIVE status"
  }
}

run "validate_deprecated_engines" {
  command = apply

  assert {
    condition = alltrue([
      for engine in data.mgc_dbaas_engines.deprecated_engines.engines : engine.status == "DEPRECATED"
    ])
    error_message = "Found engines that are not in DEPRECATED status"
  }
}

run "validate_active_engines_not_empty" {
  command = apply

  assert {
    condition     = length(data.mgc_dbaas_engines.active_engines.engines) > 0
    error_message = "No ACTIVE engines found"
  }
}

run "validate_all_engines_not_empty" {
  command = apply

  assert {
    condition     = length(data.mgc_dbaas_engines.all_engines.engines) > 0
    error_message = "No engines found"
  }
}

run "validate_all_engines_includes_both_statuses" {
  command = apply

  assert {
    condition = length([
      for engine in data.mgc_dbaas_engines.all_engines.engines : engine
      if contains(["ACTIVE", "DEPRECATED"], engine.status)
    ]) > 0
    error_message = "No engines with expected statuses found"
  }
}
