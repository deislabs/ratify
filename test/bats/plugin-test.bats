#!/usr/bin/env bats

load helpers

BATS_TESTS_DIR=${BATS_TESTS_DIR:-test/bats/tests}
WAIT_TIME=60
SLEEP_TIME=1

@test "cert rotator test" {
    helm uninstall ratify --namespace gatekeeper-system
    make e2e-helm-deploy-ratify CERT_DIR=${EXPIRING_CERT_DIR} CERT_ROTATION_ENABLED=true GATEKEEPER_VERSION=${GATEKEEPER_VERSION}
    sleep 120
    run [ "$(kubectl get secret ratify-tls -n gatekeeper-system -o json | jq '.data."ca.crt"')" != "$(cat ${EXPIRING_CERT_DIR}/ca.crt | base64 | tr -d '\n')" ]
    assert_success
}

@test "cosign test" {
    teardown() {
        echo "cleaning up"
        wait_for_process ${WAIT_TIME} ${SLEEP_TIME} 'kubectl delete pod cosign-demo-key --namespace default --force --ignore-not-found=true'
        wait_for_process ${WAIT_TIME} ${SLEEP_TIME} 'kubectl delete pod cosign-demo-unsigned --namespace default --force --ignore-not-found=true'
    }
    run kubectl apply -f ./library/default/template.yaml
    assert_success
    sleep 5
    run kubectl apply -f ./library/default/samples/constraint.yaml
    assert_success
    sleep 5

    run kubectl run cosign-demo-key --namespace default --image=registry:5000/cosign:signed-key
    assert_success

    run kubectl run cosign-demo-unsigned --namespace default --image=registry:5000/cosign:unsigned
    assert_failure
}

@test "cosign keyless test" {
    teardown() {
        echo "cleaning up"
        wait_for_process ${WAIT_TIME} ${SLEEP_TIME} 'kubectl delete pod cosign-demo-keyless --namespace default --force --ignore-not-found=true'
        wait_for_process ${WAIT_TIME} ${SLEEP_TIME} 'kubectl replace -f ./config/samples/config_v1beta1_verifier_cosign.yaml'
        wait_for_process ${WAIT_TIME} ${SLEEP_TIME} 'kubectl replace -f ./config/samples/config_v1beta1_store_oras_http.yaml'
    }

    # use imperative command to guarantee useHttp is updated
    run kubectl replace -f ./config/samples/config_v1beta1_verifier_cosign_keyless.yaml
    sleep 5

    run kubectl replace -f ./config/samples/config_v1beta1_store_oras.yaml
    sleep 5

    run kubectl run cosign-demo-keyless --namespace default --image=wabbitnetworks.azurecr.io/test/cosign-image:signed-keyless
    assert_success
}

@test "licensechecker test" {
    teardown() {
        echo "cleaning up"
        wait_for_process ${WAIT_TIME} ${SLEEP_TIME} 'kubectl delete pod license-checker --namespace default --force --ignore-not-found=true'
        wait_for_process ${WAIT_TIME} ${SLEEP_TIME} 'kubectl delete pod license-checker2 --namespace default --force --ignore-not-found=true'
        wait_for_process ${WAIT_TIME} ${SLEEP_TIME} 'kubectl delete verifiers.config.ratify.deislabs.io/verifier-license-checker --namespace default --ignore-not-found=true'
    }

    run kubectl apply -f ./library/default/template.yaml
    assert_success
    sleep 5
    run kubectl apply -f ./library/default/samples/constraint.yaml
    assert_success
    sleep 5
    run kubectl apply -f ./config/samples/config_v1beta1_verifier_partial_licensechecker.yaml
    sleep 5
    run kubectl run license-checker --namespace default --image=registry:5000/licensechecker:v0
    assert_failure

    run kubectl apply -f ./config/samples/config_v1beta1_verifier_complete_licensechecker.yaml
    # wait for the httpserver cache to be invalidated
    sleep 15
    run kubectl run license-checker2 --namespace default --image=registry:5000/licensechecker:v0
    assert_success
}

@test "sbom verifier test" {
    teardown() {
        echo "cleaning up"
        wait_for_process ${WAIT_TIME} ${SLEEP_TIME} 'kubectl delete pod sbom --namespace default --force --ignore-not-found=true'
        wait_for_process ${WAIT_TIME} ${SLEEP_TIME} 'kubectl delete pod sbom2 --namespace default --force --ignore-not-found=true'
    }

    run kubectl apply -f ./library/default/template.yaml
    assert_success
    sleep 5
    run kubectl apply -f ./library/default/samples/constraint.yaml
    assert_success
    sleep 5

    run kubectl apply -f ./config/samples/config_v1beta1_verifier_sbom.yaml
    sleep 5
    run kubectl run sbom --namespace default --image=registry:5000/sbom:v0
    assert_success

    run kubectl delete verifiers.config.ratify.deislabs.io/verifier-sbom
    assert_success
    # wait for the httpserver cache to be invalidated
    sleep 15
    run kubectl run sbom2 --namespace default --image=registry:5000/sbom:v0
    assert_failure
}

@test "schemavalidator verifier test" {
    skip "Skipping test for now until expected usage/configuration of this plugin can be verified"
    teardown() {
        echo "cleaning up"
        wait_for_process ${WAIT_TIME} ${SLEEP_TIME} 'kubectl delete verifiers.config.ratify.deislabs.io/verifier-license-checker --namespace default --ignore-not-found=true'
        wait_for_process ${WAIT_TIME} ${SLEEP_TIME} 'kubectl delete verifiers.config.ratify.deislabs.io/verifier-sbom --namespace default --ignore-not-found=true'
        wait_for_process ${WAIT_TIME} ${SLEEP_TIME} 'kubectl delete verifiers.config.ratify.deislabs.io/verifier-schemavalidator --namespace default --ignore-not-found=true'
        wait_for_process ${WAIT_TIME} ${SLEEP_TIME} 'kubectl delete pod schemavalidator --namespace default --force --ignore-not-found=true'
        wait_for_process ${WAIT_TIME} ${SLEEP_TIME} 'kubectl delete pod schemavalidator2 --namespace default --force --ignore-not-found=true'
    }

    run kubectl apply -f ./library/default/template.yaml
    assert_success
    sleep 5
    run kubectl apply -f ./library/default/samples/constraint.yaml
    assert_success
    sleep 5

    run kubectl apply -f ./config/samples/config_v1beta1_verifier_schemavalidator.yaml
    sleep 5
    run kubectl run schemavalidator --namespace default --image=registry:5000/schemavalidator:v0
    assert_success

    run kubectl apply -f ./config/samples/config_v1beta1_verifier_schemavalidator_bad.yaml
    assert_success
    # wait for the httpserver cache to be invalidated
    sleep 15
    run kubectl run schemavalidator2 --namespace default --image=registry:5000/schemavalidator:v0
    assert_failure
}

@test "sbom/notary/cosign/licensechecker/schemavalidator verifiers test" {
    teardown() {
        echo "cleaning up"
        wait_for_process ${WAIT_TIME} ${SLEEP_TIME} 'kubectl delete verifiers.config.ratify.deislabs.io/verifier-license-checker --namespace default --ignore-not-found=true'
        wait_for_process ${WAIT_TIME} ${SLEEP_TIME} 'kubectl delete verifiers.config.ratify.deislabs.io/verifier-sbom --namespace default --ignore-not-found=true'
        # Skipping test for now until expected usage/configuration of this plugin can be verified
        # wait_for_process ${WAIT_TIME} ${SLEEP_TIME} 'kubectl delete verifiers.config.ratify.deislabs.io/verifier-schemavalidator --namespace default --ignore-not-found=true'
        wait_for_process ${WAIT_TIME} ${SLEEP_TIME} 'kubectl delete pod all-in-one --namespace default --force --ignore-not-found=true'
    }

    run kubectl apply -f ./library/default/template.yaml
    assert_success
    sleep 5
    run kubectl apply -f ./library/default/samples/constraint.yaml
    assert_success
    sleep 5

    run kubectl apply -f ./config/samples/config_v1beta1_verifier_cosign.yaml
    sleep 5
    run kubectl apply -f ./config/samples/config_v1beta1_verifier_sbom.yaml
    sleep 5
    run kubectl apply -f ./config/samples/config_v1beta1_verifier_complete_licensechecker.yaml

    # Skipping test for now until expected usage/configuration of this plugin can be verified
    # sleep 5
    # run kubectl apply -f ./config/samples/config_v1beta1_verifier_schemavalidator.yaml
    # sleep 5

    # wait for the httpserver cache to be invalidated
    sleep 15
    run kubectl run all-in-one --namespace default --image=registry:5000/all:v0
    assert_success
}

@test "validate crd add, replace and delete" {
    teardown() {
        echo "cleaning up"
        wait_for_process ${WAIT_TIME} ${SLEEP_TIME} 'kubectl delete pod crdtest --namespace default --force --ignore-not-found=true'
    }

    echo "adding license checker, delete notary verifier and validate deployment fails due to missing notary verifier"
    run kubectl apply -f ./config/samples/config_v1beta1_verifier_complete_licensechecker.yaml
    assert_success
    run kubectl delete verifiers.config.ratify.deislabs.io/verifier-notary
    assert_success
    # wait for the httpserver cache to be invalidated
    sleep 15
    run kubectl run crdtest --namespace default --image=registry:5000/notation:signed
    assert_failure

    echo "Add notary verifier and validate deployment succeeds"
    run kubectl apply -f ./config/samples/config_v1beta1_verifier_notary.yaml
    assert_success

    # wait for the httpserver cache to be invalidated
    sleep 15
    run kubectl run crdtest --namespace default --image=registry:5000/notation:signed
    assert_success
}

@test "dynamic plugins disabled test" {
    teardown() {
        echo "cleaning up"
        wait_for_process ${WAIT_TIME} ${SLEEP_TIME} 'kubectl delete verifiers.config.ratify.deislabs.io/verifier-dynamic --namespace default --ignore-not-found=true'
    }

    start=$(date --iso-8601=seconds)
    latestpod=$(kubectl -n gatekeeper-system get pod -l=app.kubernetes.io/name=ratify --sort-by=.metadata.creationTimestamp -o=name | tail -n 1)

    run kubectl apply -f ./config/samples/config_v1beta1_verifier_dynamic.yaml
    sleep 5

    run bash -c "kubectl -n gatekeeper-system logs $latestpod --since-time=$start | grep 'dynamic plugins are currently disabled'"
    assert_success
}
