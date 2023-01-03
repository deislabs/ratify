#!/usr/bin/env bats

load helpers

BATS_TESTS_DIR=${BATS_TESTS_DIR:-test/bats/tests}
WAIT_TIME=60
SLEEP_TIME=1

@test "notary test" {
    run kubectl apply -f ./library/default/template.yaml
    assert_success
    sleep 5
    run kubectl apply -f ./library/default/samples/constraint.yaml
    assert_success
    sleep 5
    run kubectl run demo --namespace default --image=wabbitnetworks.azurecr.io/test/notary-image:signed
    assert_success
    run kubectl run demo1 --namespace default --image=wabbitnetworks.azurecr.io/test/notary-image:unsigned
    assert_failure

    echo "cleaning up"
    wait_for_process ${WAIT_TIME} ${SLEEP_TIME} kubectl delete pod demo --namespace default
}

@test "cosign test" {
    run kubectl apply -f ./library/default/template.yaml
    assert_success
    sleep 5
    run kubectl apply -f ./library/default/samples/constraint.yaml
    assert_success
    sleep 5
    run kubectl run cosign-demo --namespace default --image=wabbitnetworks.azurecr.io/test/cosign-image:signed
    assert_success
    run kubectl run cosign-demo2 --namespace default --image=wabbitnetworks.azurecr.io/test/cosign-image:unsigned
    assert_failure

    echo "cleaning up"
    wait_for_process ${WAIT_TIME} ${SLEEP_TIME} kubectl delete pod cosign-demo --namespace default
}

@test "licensechecker test" {
    run kubectl apply -f ./library/default/template.yaml
    assert_success
    sleep 5
    run kubectl apply -f ./library/default/samples/constraint.yaml
    assert_success
    sleep 5

    run kubectl apply -f ./config/samples/config_v1alpha1_verifier_partial_licensechecker.yaml
    sleep 5
    run kubectl run license-checker --namespace default --image=wabbitnetworks.azurecr.io/test/license-checker-image:v1
    assert_failure

    run kubectl apply -f ./config/samples/config_v1alpha1_verifier_complete_licensechecker.yaml
    sleep 5
    run kubectl run license-checker2 --namespace default --image=wabbitnetworks.azurecr.io/test/license-checker-image:v1
    assert_success

    echo "cleaning up"
    wait_for_process ${WAIT_TIME} ${SLEEP_TIME} kubectl delete pod license-checker --namespace default
    wait_for_process ${WAIT_TIME} ${SLEEP_TIME} kubectl delete verifiers.config.ratify.deislabs.io/verifier-license-checker --namespace default
}

@test "sbom verifier test" {
    run kubectl apply -f ./library/default/template.yaml
    assert_success
    sleep 5
    run kubectl apply -f ./library/default/samples/constraint.yaml
    assert_success
    sleep 5

    run kubectl apply -f ./config/samples/config_v1alpha1_verifier_sbom.yaml
    sleep 5
    run kubectl run sbom --namespace default --image=wabbitnetworks.azurecr.io/test/sbom-image:signed
    assert_success

    run kubectl delete verifiers.config.ratify.deislabs.io/verifier-sbom
    assert_success
    run kubectl run sbom2 --namespace default --image=wabbitnetworks.azurecr.io/test/sbom-image:signed
    assert_failure

    echo "cleaning up"
    wait_for_process ${WAIT_TIME} ${SLEEP_TIME} kubectl delete pod sbom --namespace default
}

@test "schemavalidator verifier test" {
    run kubectl apply -f ./library/default/template.yaml
    assert_success
    sleep 5
    run kubectl apply -f ./library/default/samples/constraint.yaml
    assert_success
    sleep 5

    run kubectl apply -f ./config/samples/config_v1alpha1_verifier_schemavalidator.yaml
    sleep 5
    run kubectl run schemavalidator --namespace default --image=wabbitnetworks.azurecr.io/test/all-in-one-image:signed
    assert_success

    run kubectl apply -f ./config/samples/config_v1alpha1_verifier_schemavalidator_bad.yaml
    sleep 5
    run kubectl run schemavalidator2 --namespace default --image=wabbitnetworks.azurecr.io/test/all-in-one-image:signed
    assert_failure

    echo "cleaning up"
    wait_for_process ${WAIT_TIME} ${SLEEP_TIME} kubectl delete pod schemavalidator --namespace default    
    wait_for_process ${WAIT_TIME} ${SLEEP_TIME} kubectl delete verifiers.config.ratify.deislabs.io/verifier-schemavalidator --namespace default
}

@test "sbom/notary/cosign/licensechecker/schemavalidator verifiers test" {
    run kubectl apply -f ./library/default/template.yaml
    assert_success
    sleep 5
    run kubectl apply -f ./library/default/samples/constraint.yaml
    assert_success
    sleep 5

    run kubectl apply -f ./config/samples/config_v1alpha1_verifier_sbom.yaml
    sleep 5
    run kubectl apply -f ./config/samples/config_v1alpha1_verifier_complete_licensechecker.yaml
    sleep 5
    run kubectl apply -f ./config/samples/config_v1alpha1_verifier_schemavalidator.yaml
    sleep 5

    run kubectl run all-in-one --namespace default --image=wabbitnetworks.azurecr.io/test/all-in-one-image:signed
    assert_success

    echo "cleaning up"
    wait_for_process ${WAIT_TIME} ${SLEEP_TIME} kubectl delete verifiers.config.ratify.deislabs.io/verifier-license-checker --namespace default
    wait_for_process ${WAIT_TIME} ${SLEEP_TIME} kubectl delete verifiers.config.ratify.deislabs.io/verifier-sbom --namespace default
    wait_for_process ${WAIT_TIME} ${SLEEP_TIME} kubectl delete verifiers.config.ratify.deislabs.io/verifier-schemavalidator --namespace default
    wait_for_process ${WAIT_TIME} ${SLEEP_TIME} kubectl delete pod all-in-one --namespace default
}

@test "validate crd add, replace and delete" {
    echo "adding license checker, delete notary verifier and validate deployment fails due to missing notary verifier"
    run kubectl apply -f ./config/samples/config_v1alpha1_verifier_complete_licensechecker.yaml
    assert_success
    run kubectl delete verifiers.config.ratify.deislabs.io/verifier-notary
    assert_success
    run kubectl run crdtest --namespace default --image=wabbitnetworks.azurecr.io/test/notary-image:signed
    assert_failure

    echo "Add notary verifier and validate deployment succeeds"
    run kubectl apply -f ./config/samples/config_v1alpha1_verifier_notary.yaml
    assert_success

    run kubectl run crdtest --namespace default --image=wabbitnetworks.azurecr.io/test/notary-image:signed
    assert_success

    echo "cleaning up"
    wait_for_process ${WAIT_TIME} ${SLEEP_TIME} kubectl delete pod crdtest --namespace default
}

@test "configmap update test" {
    skip "Skipping test for now as we are no longer watching for configfile update in a k8 environment.This test ensures we are watching config file updates in a non-kub scenario"
    run kubectl apply -f ./library/default/template.yaml
    assert_success
    sleep 5
    run kubectl apply -f ./library/default/samples/constraint.yaml
    assert_success
    sleep 5
    run kubectl run demo2 --image=wabbitnetworks.azurecr.io/test/net-monitor:signed
    assert_success

    run kubectl get configmaps ratify-configuration --namespace=ratify-service -o yaml >currentConfig.yaml
    run kubectl delete -f ./library/default/samples/constraint.yaml

    wait_for_process ${WAIT_TIME} ${SLEEP_TIME} "kubectl replace --namespace=ratify-service -f ${BATS_TESTS_DIR}/configmap/invalidconfigmap.yaml"
    echo "Waiting for 150 second for configuration update"
    sleep 150

    run kubectl apply -f ./library/default/samples/constraint.yaml
    assert_success
    run kubectl run demo3 --image=wabbitnetworks.azurecr.io/test/net-monitor:signed
    echo "Current time after validate : $(date +"%T")"
    assert_failure

    wait_for_process ${WAIT_TIME} ${SLEEP_TIME} "kubectl replace --namespace=ratify-service -f currentConfig.yaml"
}
