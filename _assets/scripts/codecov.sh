#!/usr/bin/env bash

source "${GIT_ROOT}/_assets/scripts/colors.sh"

report_to_codecov() {
  # https://go.dev/blog/integration-test-coverage
  echo -e "${GRN}Uploading coverage report to Codecov${RST}"

  local tests_report_wildcard="${1}"
  local coverage_report="${2}"
  local flag="${3}"

  # Gather report files with given wildcard
  local report_files_args=""
  for file in ${tests_report_wildcard}; do
    report_files_args+="--file ${file} "
  done

  codecov do-upload --token "${CODECOV_TOKEN}" --report-type test_results ${report_files_args}
  codecov upload-process --token "${CODECOV_TOKEN}" -f ${coverage_report} -F "${flag}"
}

convert_coverage_to_html() {
  echo -e "${GRN}Generating HTML coverage report${RST}"

  local input_coverage_report="${1}"
  local output_coverage_report="${2}"

  go tool cover -html "${input_coverage_report}" -o "${output_coverage_report}"
}