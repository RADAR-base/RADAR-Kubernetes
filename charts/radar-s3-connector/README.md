

# radar-s3-connector

![Version: 0.1.1](https://img.shields.io/badge/Version-0.1.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 1.0](https://img.shields.io/badge/AppVersion-1.0-informational?style=flat-square)

A Helm chart for RADAR-base s3 connector. This connector uses Confluent s3 connector with a custom data transformers. These configurations enable a sink connector. See full list of properties here https://docs.confluent.io/kafka-connect-s3-sink/current/configuration_options.html#s3-configuration-options

**Homepage:** <https://radar-base.org>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Keyvan Hedayati | keyvan@thehyve.nl |  |
| Joris Borgdorff | joris@thehyve.nl |  |
| Nivethika Mahasivam | nivethika@thehyve.nl |  |

## Source Code

* <https://github.com/RADAR-base/kafka-connect-transform-keyvalue>
* <https://docs.confluent.io/kafka-connect-s3-sink/current/configuration_options.html#s3-configuration-options>

## Prerequisites
* Kubernetes 1.17+
* Kubectl 1.17+
* Helm 3.1.0+
* PV provisioner support in the underlying infrastructure

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| replicaCount | int | `1` |  |
| image.repository | string | `"radarbase/kafka-connect-transform-s3"` |  |
| image.tag | string | `"5.5.1"` |  |
| image.pullPolicy | string | `"IfNotPresent"` |  |
| nameOverride | string | `""` |  |
| fullnameOverride | string | `""` |  |
| service.type | string | `"ClusterIP"` |  |
| service.port | int | `8083` |  |
| ingress.enabled | bool | `false` |  |
| ingress.annotations | object | `{}` |  |
| ingress.hosts[0].host | string | `"chart-example.local"` |  |
| ingress.hosts[0].paths | list | `[]` |  |
| ingress.tls | list | `[]` |  |
| resources.requests.cpu | string | `"100m"` |  |
| resources.requests.memory | string | `"3Gi"` |  |
| nodeSelector | object | `{}` |  |
| tolerations | list | `[]` |  |
| affinity | object | `{}` |  |
| kafka.url | string | `"PLAINTEXT://cp-kafka-headless:9092"` | Kafka broker URLs |
| kafka.securityProtocol | string | `"PLAINTEXT"` | Not used. To be confirmed |
| kafka.saslJaasConfig | string | `""` | Not used. To be confirmed |
| kafka.saslMechanism | string | `"GSSAPI"` | Not used. To be confirmed |
| kafka.sslEndpointIdentificationAlgorithm | string | `"https"` | Not used. To be confirmed |
| schemaRegistry.url | string | `"http://cp-schema-registry:8081"` |  |
| schemaRegistry.basicAuth | string | `""` | Not used. To be confirmed |
| topics | string | `"android_phone_usage_event_output,android_biovotion_vsm1_acceleration,android_biovotion_vsm1_battery_level,android_biovotion_vsm1_blood_volume_pulse,android_biovotion_vsm1_energy,android_biovotion_vsm1_galvanic_skin_response,android_biovotion_vsm1_heartrate,android_biovotion_vsm1_heartrate_variability,android_biovotion_vsm1_led_current,android_biovotion_vsm1_oxygen_saturation,android_biovotion_vsm1_ppg_raw,android_biovotion_vsm1_respiration_rate,android_biovotion_vsm1_temperature,android_bittium_faros_acceleration,android_bittium_faros_battery_level,android_bittium_faros_ecg,android_bittium_faros_inter_beat_interval,android_bittium_faros_temperature,android_empatica_e4_acceleration,android_empatica_e4_battery_level,android_empatica_e4_blood_volume_pulse,android_empatica_e4_electrodermal_activity,android_empatica_e4_inter_beat_interval,android_empatica_e4_sensor_status,android_empatica_e4_temperature,android_local_weather,android_pebble_2_acceleration,android_pebble_2_battery_level,android_pebble_2_heartrate,android_pebble_2_heartrate_filtered,android_phone_acceleration,android_phone_battery_level,android_phone_bluetooth_devices,android_phone_call,android_phone_contacts,android_phone_gyroscope,android_phone_light,android_phone_magnetic_field,android_phone_ppg,android_phone_relative_location,android_phone_sms,android_phone_sms_unread,android_phone_step_count,android_phone_usage_event,android_phone_user_interaction,android_processed_audio,application_device_info,application_external_time,application_record_counts,application_server_status,application_time_zone,application_uptime,certh_banking_app_event,certh_banking_app_transaction,connect_fitbit_activity_log,connect_fitbit_intraday_calories,connect_fitbit_intraday_heart_rate,connect_fitbit_intraday_steps,connect_fitbit_sleep_classic,connect_fitbit_sleep_stages,connect_fitbit_time_zone,connect_upload_altoida_acceleration,connect_upload_altoida_action,connect_upload_altoida_attitude,connect_upload_altoida_bit_metrics,connect_upload_altoida_blink,connect_upload_altoida_diagnostics,connect_upload_altoida_domain_result,connect_upload_altoida_dot_metrics,connect_upload_altoida_eye_tracking,connect_upload_altoida_gravity,connect_upload_altoida_magnetic_field,connect_upload_altoida_metadata,connect_upload_altoida_object,connect_upload_altoida_path,connect_upload_altoida_rotation,connect_upload_altoida_summary,connect_upload_altoida_tap,connect_upload_altoida_touch,connect_upload_axivity_acceleration,connect_upload_axivity_battery_level,connect_upload_axivity_event,connect_upload_axivity_light,connect_upload_axivity_metadata,connect_upload_axivity_temperature,connect_upload_oxford_camera_data,connect_upload_oxford_camera_image,connect_upload_physilog_binary_data,notification_thinc_it,questionnaire_app_event,questionnaire_ari_self,questionnaire_art_cognitive_test,questionnaire_audio,questionnaire_baars_iv,questionnaire_bipq,questionnaire_completion_log,questionnaire_esm,questionnaire_esm28q,questionnaire_esm_epi_mod_1,questionnaire_evening_assessment,questionnaire_gad7,questionnaire_morning_assessment,questionnaire_patient_determined_disease_step,questionnaire_perceived_deficits_questionnaire,questionnaire_phq8,questionnaire_rpq,questionnaire_rses,questionnaire_tam,questionnaire_timezone,task_2MW_test,task_romberg_test,task_tandem_walking_test,thincit_code_breaker,thincit_pdq5,thincit_spotter,thincit_symbol_check,thincit_trails,"` | List of topics to be consumed by the sink connector separated by comma. |
| s3Endpoint | string | `"http://minio:9000/"` | Target S3 endpoint url |
| s3Tagging | bool | `false` | set to true, if S3 objects should be tagged with start and end offsets, as well as record count. |
| s3PartSize | int | `5242880` | The Part Size in S3 Multi-part Uploads. |
| s3Region | string | `nil` | The AWS region to be used the connector. |
| flushSize | int | `10000` | Number of records written to store before invoking file commits. |
| rotateInterval | int | `900000` | The time interval in milliseconds to invoke file commits. |
| maxTasks | int | `4` | Number of tasks in the connector |
| bucketAccessKey | string | `"access_key"` | Access key of the target S3 bucket |
| bucketSecretKey | string | `"secret"` | Secret key of the target S3 bucket |
| bucketName | string | `"radar_intermediate_storage"` | Bucket name of the target S3 bucket |
| cc.enabled | bool | `false` | Set to true, if Confluent Cloud is used |
| cc.bootstrapServerurl | string | `""` | Confluent cloud based Kafka broker URL (if Confluent Cloud based Kafka cluster is used) |
| cc.schemaRegistryUrl | string | `""` | Confluent cloud based Schema registry URL (if Confluent Cloud based Schema registry is used) |
| cc.apiKey | string | `"ccApikey"` | API Key of the Confluent Cloud cluster |
| cc.apiSecret | string | `"ccApiSecret"` | API secret of the Confluent Cloud cluster |
| cc.schemaRegistryApiKey | string | `"srApiKey"` | API Key of the Confluent Cloud Schema registry |
| cc.schemaRegistryApiSecret | string | `"srApiSecret"` | API Key of the Confluent Cloud Schema registry |
