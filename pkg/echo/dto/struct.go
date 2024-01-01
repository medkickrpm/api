package dto

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CareplanSheetResponse struct {
	Submission_Date                                       string `json:"submission_date"`
	Office_Email                                          string `json:"office_email"`
	Care_Coordinator                                      string `json:"care_coordinator"`
	First_Name                                            string `json:"first_name"`
	Last_Name                                             string `json:"last_name"`
	DOB                                                   string `json:"dob"`
	Chronic_Conditions                                    string `json:"chronic_conditions"`
	Other_Diagnosis_Not_Listed                            string `json:"other_diagnosis_not_listed"`
	No_Label                                              string `json:"no_label"`
	Goals_Areas                                           string `json:"goals_areas"`
	Medication_Reconciliation                             string `json:"medication_reconciliation"`
	Activities_of_Daily_Living_Improvement_Goals          string `json:"activities_of_daily_living_improvement_goals"`
	Obstacles                                             string `json:"obstacles"`
	Food_Preparation_Ideas                                string `json:"food_preparation_ideas"`
	ADL_Nutrition_Plan                                    string `json:"adl_nutrition_plan"`
	Bathing_Showering                                     string `json:"bathing_showering"`
	Toileting                                             string `json:"toileting"`
	Dressing                                              string `json:"dressing"`
	Ambulation_walking_and_using_stairs                   string `json:"ambulation_walking_and_using_stairs"`
	Standing_Getting_out_of_chair_Balance                 string `json:"standing_getting_out_of_chair_balance"`
	Medication_Management                                 string `json:"medication_management"`
	Cleaning                                              string `json:"cleaning"`
	Running_Errands_Transportation                        string `json:"running_errands_transportation"`
	Using_Phone_Technology                                string `json:"using_phone_technology"`
	Managing_Finances                                     string `json:"managing_finances"`
	Social_Activities                                     string `json:"social_activities"`
	Exercise                                              string `json:"exercise"`
	Hobbies                                               string `json:"hobbies"`
	Pet_Care                                              string `json:"pet_care"`
	Post_Heart_Attack_Goals                               string `json:"post_heart_attack_goals"`
	Post_Heart_Attack_Plan                                string `json:"post_heart_attack_plan"`
	CHF_Goals                                             string `json:"chf_goals"`
	CHF_Improvement_Plan                                  string `json:"chf_improvement_plan"`
	Coronary_Artery_Disease_Goals                         string `json:"coronary_artery_disease_goals"`
	Coronary_Artery_Disease_Plan                          string `json:"coronary_artery_disease_plan"`
	Kidney_Health_Goals                                   string `json:"kidney_health_goals"`
	Kidney_Health_Plan                                    string `json:"kidney_health_plan"`
	Hypertension_Goals                                    string `json:"hypertension_goals"`
	Hypertension_Plan                                     string `json:"hypertension_plan"`
	Diabetes_Goals                                        string `json:"diabetes_goals"`
	Diabetes_Improvement_Plan                             string `json:"diabetes_improvement_plan"`
	Respiratory_Goals                                     string `json:"respiratory_goals"`
	Respiratory_Challenges_and_Breathing_Improvement_Plan string `json:"respiratory_challenges_and_breathing_improvement_plan"`
	Oxygen_Safety_Goals                                   string `json:"oxygen_safety_goals"`
	Oxygen_Safety_Plan                                    string `json:"oxygen_safety_plan"`
	GERD_Goals                                            string `json:"gerd_goals"`
	GERD_Plan                                             string `json:"gerd_plan"`
	Cholesterol_Goals                                     string `json:"cholesterol_goals"`
	Cholesterol_Improvement_Plan                          string `json:"cholesterol_improvement_plan"`
	Self_Care_Goals                                       string `json:"self_care_goals"`
	Self_Care_Plan                                        string `json:"self_care_plan"`
	Stress_Goals                                          string `json:"stress_goals"`
	Stress_Reduction_Plan                                 string `json:"stress_reduction_plan"`
	Fall_Prevention_Goals                                 string `json:"fall_prevention_goals"`
	Fall_Prevention_Plan                                  string `json:"fall_prevention_plan"`
	Sleep_Goals                                           string `json:"sleep_goals"`
	Sleep_Improvement_Plan                                string `json:"sleep_improvement_plan"`
	Weight_Loss_Goals                                     string `json:"weight_loss_goals"`
	Weight_Loss_Plan                                      string `json:"weight_loss_plan"`
	Diet_Goals                                            string `json:"diet_goals"`
	Diet_Improvement_Plan                                 string `json:"diet_improvement_plan"`
	Mental_Health_Goals                                   string `json:"mental_health_goals"`
	Mental_Health_Improvement_Plan                        string `json:"mental_health_improvement_plan"`
	Low_Blood_Pressure_Improvement_Goals                  string `json:"low_blood_pressure_improvement_goals"`
	Low_Blood_Pressure_Avoidance_Plan                     string `json:"low_blood_pressure_avoidance_plan"`
	Anemia_Goals                                          string `json:"anemia_goals"`
	Anemia_Improvement_Plan                               string `json:"anemia_improvement_plan"`
	Thyroid_Health_Goals                                  string `json:"thyroid_health_goals"`
	Memory_Care_Goals                                     string `json:"memory_care_goals"`
	Memory_Care_Plan                                      string `json:"memory_care_plan"`
	A_Fib_Goals                                           string `json:"a-fib_goals"`
	A_Fib_Improvement_Plan                                string `json:"a-fib_improvement_plan"`
	Eye_Health_Goals                                      string `json:"eye_health_goals"`
	Eye_Health_Improvement_Plan                           string `json:"eye_health_improvement_plan"`
	Bone_Health_Goals                                     string `json:"bone_health_goals"`
	Bone_Health_Plan                                      string `json:"bone_health_plan"`
	Post_Stroke_TIA_goals                                 string `json:"post_stroke_tia_goals"`
	Post_Stroke_TIA_Plan                                  string `json:"post_stroke_tia_plan"`
	Chronic_Headache_Migraine_Goals                       string `json:"chronic_headache_migraine_goals"`
	Chronic_Headache_Migraine_Improvement_Plan            string `json:"chronic_headache_migraine_improvement_plan"`
	Thyroid_Health_Improvement_Plan                       string `json:"thyroid_health_improvement_plan"`
	Pain_Goals                                            string `json:"pain_goals"`
	Pain_Reduction_Plan                                   string `json:"pain_reduction_plan"`
	Exercise_Goals                                        string `json:"exercise_goals"`
	Exercise_Activity_Plan                                string `json:"exercise_activity_plan"`
	Energy_Improvement_Plan                               string `json:"energy_improvement_plan"`
	Bowel_Regularity_Plan                                 string `json:"bowel_regularity_plan"`
	Improving_Hair_and_Skin_Plan                          string `json:"improving_hair_and_skin_plan"`
	Hydration_goals                                       string `json:"hydration_goals"`
	Hydration_Improvement_Plan                            string `json:"hydration_improvement_plan"`
	Sodium_Reduction_Plan                                 string `json:"sodium_reduction_plan"`
	Date                                                  string `json:"date"`
	Tobacco_Use_Goals                                     string `json:"tobacco_use_goals"`
	Tobacco_Use_Plan                                      string `json:"tobacco_use_plan"`
	Input153_Goals                                        string `json:"input153_goals"`
	Input153_Plan                                         string `json:"input153_plan"`
	No_Label_2                                            string `json:"no_label_2"`
	Preventatives_Due                                     string `json:"preventatives_due"`
	Submission_ID                                         string `json:"submission_id"`
}

type VerifyUserFieldResponse struct {
	IsAvailable bool `json:"is_available"`
}
