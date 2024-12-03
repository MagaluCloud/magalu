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

run "validate_active_and_deprecated_instance_types_are_different" {
  command = apply

  assert {
    condition = length(setintersection(
      toset([for it in data.mgc_dbaas_instance_types.active_instance_types.instance_types : it.id]),
      toset([for it in data.mgc_dbaas_instance_types.deprecated_instance_types.instance_types : it.id])
    )) == 0
    error_message = "Found instance types that appear in both active and deprecated lists"
  }
}

run "validate_active_instance_types_not_empty" {
  command = apply

  assert {
    condition     = length(data.mgc_dbaas_instance_types.active_instance_types.instance_types) > 0
    error_message = "No active instance types found"
  }
}

run "validate_default_instance_types_matches_active" {
  command = apply

  assert {
    condition = length(setintersection(
      toset([for it in data.mgc_dbaas_instance_types.default_instance_types.instance_types : it.id]),
      toset([for it in data.mgc_dbaas_instance_types.active_instance_types.instance_types : it.id])
    )) == length(data.mgc_dbaas_instance_types.default_instance_types.instance_types)
    error_message = "Default instance types list does not match active instance types"
  }
}

run "validate_all_instance_types_not_empty" {
  command = apply

  assert {
    condition     = length(data.mgc_dbaas_instance_types.default_instance_types.instance_types) > 0
    error_message = "No instance types found"
  }
}


run "validate_instance_type_fields" {
  command = apply

  assert {
    condition = alltrue([
      for it in data.mgc_dbaas_instance_types.default_instance_types.instance_types :
      it.id != "" && it.name != "" && it.ram != "" && it.vcpu != "" && it.size != ""
    ])
    error_message = "Found instance type with empty required fields"
  }
}
