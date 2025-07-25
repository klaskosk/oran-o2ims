/*
SPDX-FileCopyrightText: Red Hat

SPDX-License-Identifier: Apache-2.0
*/

package controllers

import (
	"context"

	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	hwmgmtv1alpha1 "github.com/openshift-kni/oran-o2ims/api/hardwaremanagement/v1alpha1"
	"github.com/openshift-kni/oran-o2ims/internal/controllers/utils"
	testutils "github.com/openshift-kni/oran-o2ims/test/utils"
)

/*
import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	hwmgmtv1alpha1 "github.com/openshift-kni/oran-o2ims/api/hardwaremanagement/v1alpha1"
	provisioningv1alpha1 "github.com/openshift-kni/oran-o2ims/api/provisioning/v1alpha1"
	hwmgrpluginapi "github.com/openshift-kni/oran-o2ims/hwmgr-plugins/api/generated/client"
	"github.com/openshift-kni/oran-o2ims/internal/controllers/utils"
	testutils "github.com/openshift-kni/oran-o2ims/test/utils"
	"github.com/openshift/assisted-service/api/v1beta1"
	siteconfig "github.com/stolostron/siteconfig/api/v1alpha1"
)

const (
	groupNameController = "controller"
	groupNameWorker     = "worker"
)


var _ = Describe("renderHardwareTemplate", func() {
	var (
		ctx             context.Context
		c               client.Client
		reconciler      *ProvisioningRequestReconciler
		task            *provisioningRequestReconcilerTask
		clusterInstance *siteconfig.ClusterInstance
		ct              *provisioningv1alpha1.ClusterTemplate
		cr              *provisioningv1alpha1.ProvisioningRequest
		tName           = "clustertemplate-a"
		tVersion        = "v1.0.0"
		ctNamespace     = "clustertemplate-a-v4-16"
		hwTemplate      = "hwTemplate-v1"
		hwTemplatev2    = "hwTemplate-v2"
		crName          = "cluster-1"
	)

	BeforeEach(func() {
		ctx = context.Background()

		// Define the cluster instance.
		clusterInstance = &siteconfig.ClusterInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      crName,
				Namespace: ctNamespace,
			},
			Spec: siteconfig.ClusterInstanceSpec{
				Nodes: []siteconfig.NodeSpec{
					{
						Role: "master",
						NodeNetwork: &v1beta1.NMStateConfigSpec{
							Interfaces: []*v1beta1.Interface{
								{Name: "eno12399"},
								{Name: "eno1"},
							},
						},
					},
					{
						Role: "master",
						NodeNetwork: &v1beta1.NMStateConfigSpec{
							Interfaces: []*v1beta1.Interface{
								{Name: "eno12399"},
								{Name: "eno1"},
							},
						},
					},
					{
						Role: "worker",
						NodeNetwork: &v1beta1.NMStateConfigSpec{
							Interfaces: []*v1beta1.Interface{
								{Name: "eno1"},
							},
						},
					},
				},
			},
		}

		// Define the provisioning request.
		cr = &provisioningv1alpha1.ProvisioningRequest{
			ObjectMeta: metav1.ObjectMeta{
				Name: crName,
			},
			Spec: provisioningv1alpha1.ProvisioningRequestSpec{
				TemplateName:    tName,
				TemplateVersion: tVersion,
				TemplateParameters: runtime.RawExtension{
					Raw: []byte(testutils.TestFullTemplateParameters),
				},
			},
		}

		// Define the cluster template.
		ct = &provisioningv1alpha1.ClusterTemplate{
			ObjectMeta: metav1.ObjectMeta{
				Name:      GetClusterTemplateRefName(tName, tVersion),
				Namespace: ctNamespace,
			},
			Spec: provisioningv1alpha1.ClusterTemplateSpec{
				Name:    tName,
				Version: tVersion,
				Templates: provisioningv1alpha1.Templates{
					HwTemplate: hwTemplate,
				},
			},
			Status: provisioningv1alpha1.ClusterTemplateStatus{
				Conditions: []metav1.Condition{
					{
						Type:   string(provisioningv1alpha1.CTconditionTypes.Validated),
						Reason: string(provisioningv1alpha1.CTconditionReasons.Completed),
						Status: metav1.ConditionTrue,
					},
				},
			},
		}

		c = getFakeClientFromObjects([]client.Object{cr}...)
		reconciler = &ProvisioningRequestReconciler{
			Client: c,
			Logger: logger,
		}
		task = &provisioningRequestReconcilerTask{
			logger: reconciler.Logger,
			client: reconciler.Client,
			object: cr,
		}
	})

	It("returns no error when renderHardwareTemplate succeeds", func() {
		// Ensure the ClusterTemplate is created
		Expect(c.Create(ctx, ct)).To(Succeed())

		// Define the hardware template resource
		hwTemplate := &hwmgmtv1alpha1.HardwareTemplate{
			ObjectMeta: metav1.ObjectMeta{
				Name:      hwTemplate,
				Namespace: utils.InventoryNamespace,
			},
			Spec: hwmgmtv1alpha1.HardwareTemplateSpec{
				HardwarePluginRef:  utils.UnitTestHwPluginRef,
				BootInterfaceLabel: "bootable-interface",
				NodeGroupData: []hwmgmtv1alpha1.NodeGroupData{
					{
						Name:           "controller",
						Role:           "master",
						ResourcePoolId: "xyz",
						HwProfile:      "profile-spr-single-processor-64G",
					},
					{
						Name:           "worker",
						Role:           "worker",
						ResourcePoolId: "xyz",
						HwProfile:      "profile-spr-dual-processor-128G",
					},
				},
				Extensions: map[string]string{
					"resourceTypeId": "ResourceGroup~2.1.1",
				},
			},
		}

		Expect(c.Create(ctx, hwTemplate)).To(Succeed())
		unstructuredCi, err := utils.ConvertToUnstructured(*clusterInstance)
		Expect(err).ToNot(HaveOccurred())
		nodeAllocationRequest, err := task.renderHardwareTemplate(ctx, unstructuredCi)
		Expect(err).ToNot(HaveOccurred())

		VerifyHardwareTemplateStatus(ctx, c, hwTemplate.Name, metav1.Condition{
			Type:    string(hwmgmtv1alpha1.Validation),
			Status:  metav1.ConditionTrue,
			Reason:  string(hwmgmtv1alpha1.Completed),
			Message: "Validated",
		})

		Expect(nodeAllocationRequest).ToNot(BeNil())

		roleCounts := make(map[string]int)
		for _, node := range clusterInstance.Spec.Nodes {
			// Count the nodes per group
			roleCounts[node.Role]++
		}
		Expect(nodeAllocationRequest.NodeGroup).To(HaveLen(2))
		expectedNodeGroups := map[string]struct {
			size int
		}{
			groupNameController: {size: roleCounts["master"]},
			groupNameWorker:     {size: roleCounts["worker"]},
		}

		for _, group := range nodeAllocationRequest.NodeGroup {
			expected, found := expectedNodeGroups[group.NodeGroupData.Name]
			Expect(found).To(BeTrue())
			Expect(group.NodeGroupData.Size).To(Equal(expected.size))
		}
	})

	It("returns an error when the HwTemplate is not found", func() {
		// Ensure the ClusterTemplate is created
		Expect(c.Create(ctx, ct)).To(Succeed())
		unstructuredCi, err := utils.ConvertToUnstructured(*clusterInstance)
		Expect(err).ToNot(HaveOccurred())
		nodeAllocationRequest, err := task.renderHardwareTemplate(ctx, unstructuredCi)
		Expect(err).To(HaveOccurred())
		Expect(nodeAllocationRequest).To(BeNil())
		Expect(err.Error()).To(ContainSubstring("failed to get the HardwareTemplate %s resource", hwTemplate))
	})

	It("returns an error when the ClusterTemplate is not found", func() {
		unstructuredCi, err := utils.ConvertToUnstructured(*clusterInstance)
		Expect(err).ToNot(HaveOccurred())
		nodeAllocationRequest, err := task.renderHardwareTemplate(ctx, unstructuredCi)
		Expect(err).To(HaveOccurred())
		Expect(nodeAllocationRequest).To(BeNil())
		Expect(err.Error()).To(ContainSubstring("failed to get the ClusterTemplate"))
	})

	Context("When NodeAllocationRequest has been created", func() {
		var nodeAllocationRequest *pluginsv1alpha1.NodeAllocationRequest

		BeforeEach(func() {
			// Create NodeAllocationRequest resource
			nodeAllocationRequest = &pluginsv1alpha1.NodeAllocationRequest{}

			nodeAllocationRequest.SetName(crName)
			nodeAllocationRequest.SetNamespace("hwmgr")
			nodeAllocationRequest.Spec.HardwarePluginRef = utils.UnitTestHwPluginRef
			nodeAllocationRequest.Annotations = map[string]string{hwmgmtv1alpha1.BootInterfaceLabelAnnotation: "bootable-interface"}
			nodeAllocationRequest.Spec.NodeGroup = []pluginsv1alpha1.NodeGroup{
				{
					NodeGroupData: hwmgmtv1alpha1.NodeGroupData{
						Name: groupNameController, HwProfile: "profile-spr-single-processor-64G",
					}, Size: 1,
				},
			}
			nodeAllocationRequest.Status.Conditions = []metav1.Condition{
				{Type: string(hwmgmtv1alpha1.Provisioned), Status: metav1.ConditionFalse, Reason: string(hwmgmtv1alpha1.InProgress)},
			}
			Expect(c.Create(ctx, nodeAllocationRequest)).To(Succeed())
		})
		It("returns an error when the hardware template contains a change in hardwarePluginRef", func() {
			ct.Spec.Templates.HwTemplate = hwTemplatev2
			// Ensure the ClusterTemplate is created
			Expect(c.Create(ctx, ct)).To(Succeed())
			// Define the new version of hardware template resource
			hwTemplate2 := &hwmgmtv1alpha1.HardwareTemplate{
				ObjectMeta: metav1.ObjectMeta{
					Name:      hwTemplatev2,
					Namespace: utils.InventoryNamespace,
				},
				Spec: hwmgmtv1alpha1.HardwareTemplateSpec{
					HardwarePluginRef:  "new id",
					BootInterfaceLabel: "bootable-interface",
					NodeGroupData: []hwmgmtv1alpha1.NodeGroupData{
						{
							Name:      "worker",
							Role:      "worker",
							HwProfile: "profile-spr-single-processor-64G",
						},
					},
					Extensions: map[string]string{
						"esourceTypeId": "ResourceGroup~2.1.1",
					},
				},
			}
			Expect(c.Create(ctx, hwTemplate2)).To(Succeed())
			unstructuredCi, err := utils.ConvertToUnstructured(*clusterInstance)
			Expect(err).ToNot(HaveOccurred())
			_, err = task.renderHardwareTemplate(ctx, unstructuredCi)
			Expect(err).To(HaveOccurred())

			VerifyHardwareTemplateStatus(ctx, c, hwTemplate2.Name, metav1.Condition{
				Type:    string(hwmgmtv1alpha1.Validation),
				Status:  metav1.ConditionFalse,
				Reason:  string(hwmgmtv1alpha1.Failed),
				Message: "unallowed change detected",
			})

			cond := meta.FindStatusCondition(cr.Status.Conditions, string(provisioningv1alpha1.PRconditionTypes.HardwareTemplateRendered))
			Expect(cond).ToNot(BeNil())
			testutils.VerifyStatusCondition(*cond, metav1.Condition{
				Type:    string(provisioningv1alpha1.PRconditionTypes.HardwareTemplateRendered),
				Status:  metav1.ConditionFalse,
				Reason:  string(provisioningv1alpha1.CRconditionReasons.Failed),
				Message: "Failed to render the Hardware template",
			})
		})

		It("returns an error when the hardware template contains a change in bootIntefaceLabel", func() {
			ct.Spec.Templates.HwTemplate = hwTemplatev2
			// Ensure the ClusterTemplate is created
			Expect(c.Create(ctx, ct)).To(Succeed())

			// Define the new version of hardware template resource
			hwTemplate2 := &hwmgmtv1alpha1.HardwareTemplate{
				ObjectMeta: metav1.ObjectMeta{
					Name:      hwTemplatev2,
					Namespace: utils.InventoryNamespace,
				},
				Spec: hwmgmtv1alpha1.HardwareTemplateSpec{
					HardwarePluginRef:  utils.UnitTestHwPluginRef,
					BootInterfaceLabel: "new-label",
					NodeGroupData: []hwmgmtv1alpha1.NodeGroupData{
						{
							Name:           "contoller",
							Role:           "master",
							ResourcePoolId: "xyz",
							HwProfile:      "profile-spr-single-processor-64G",
						},
					},
					Extensions: map[string]string{
						"resourceTypeId": "ResourceGroup~2.1.1",
					},
				},
			}
			Expect(c.Create(ctx, hwTemplate2)).To(Succeed())
			unstructuredCi, err := utils.ConvertToUnstructured(*clusterInstance)
			Expect(err).ToNot(HaveOccurred())
			_, err = task.renderHardwareTemplate(ctx, unstructuredCi)
			Expect(err).To(HaveOccurred())

			VerifyHardwareTemplateStatus(ctx, c, hwTemplate2.Name, metav1.Condition{
				Type:    string(hwmgmtv1alpha1.Validation),
				Status:  metav1.ConditionFalse,
				Reason:  string(hwmgmtv1alpha1.Failed),
				Message: "unallowed change detected",
			})

			cond := meta.FindStatusCondition(cr.Status.Conditions, string(provisioningv1alpha1.PRconditionTypes.HardwareTemplateRendered))
			Expect(cond).ToNot(BeNil())
			testutils.VerifyStatusCondition(*cond, metav1.Condition{
				Type:    string(provisioningv1alpha1.PRconditionTypes.HardwareTemplateRendered),
				Status:  metav1.ConditionFalse,
				Reason:  string(provisioningv1alpha1.CRconditionReasons.Failed),
				Message: "Failed to render the Hardware template",
			})
		})

		It("returns an error when the hardware template contains a change in groups", func() {
			ct.Spec.Templates.HwTemplate = hwTemplatev2
			// Ensure the ClusterTemplate is created
			Expect(c.Create(ctx, ct)).To(Succeed())

			// Define the new version of hardware template resource
			hwTemplate2 := &hwmgmtv1alpha1.HardwareTemplate{
				ObjectMeta: metav1.ObjectMeta{
					Name:      hwTemplatev2,
					Namespace: utils.InventoryNamespace,
				},
				Spec: hwmgmtv1alpha1.HardwareTemplateSpec{
					HardwarePluginRef:  utils.UnitTestHwPluginRef,
					BootInterfaceLabel: "bootable-interface",
					NodeGroupData: []hwmgmtv1alpha1.NodeGroupData{
						{
							Name:           "master",
							Role:           "master",
							ResourcePoolId: "xyz",
							HwProfile:      "profile-spr-single-processor-64G",
						},
						{
							Name:           "worker",
							Role:           "worker",
							ResourcePoolId: "xyz",
							HwProfile:      "profile-spr-single-processor-64G",
						},
					},
					Extensions: map[string]string{
						"esourceTypeId": "ResourceGroup~2.1.1",
					},
				},
			}
			Expect(c.Create(ctx, hwTemplate2)).To(Succeed())
			unstructuredCi, err := utils.ConvertToUnstructured(*clusterInstance)
			Expect(err).ToNot(HaveOccurred())
			_, err = task.renderHardwareTemplate(ctx, unstructuredCi)
			Expect(err).To(HaveOccurred())

			errMessage := fmt.Sprintf("node group %s found in NodeAllocationRequest spec but not in Hardware Template", groupNameController)
			VerifyHardwareTemplateStatus(ctx, c, hwTemplate2.Name, metav1.Condition{
				Type:    string(hwmgmtv1alpha1.Validation),
				Status:  metav1.ConditionFalse,
				Reason:  string(hwmgmtv1alpha1.Failed),
				Message: errMessage,
			})

			cond := meta.FindStatusCondition(cr.Status.Conditions, string(provisioningv1alpha1.PRconditionTypes.HardwareTemplateRendered))
			Expect(cond).ToNot(BeNil())
			testutils.VerifyStatusCondition(*cond, metav1.Condition{
				Type:    string(provisioningv1alpha1.PRconditionTypes.HardwareTemplateRendered),
				Status:  metav1.ConditionFalse,
				Reason:  string(provisioningv1alpha1.CRconditionReasons.Failed),
				Message: "Failed to render the Hardware template",
			})
		})
	})
})

var _ = Describe("waitForNodeAllocationRequestProvision", func() {
	var (
		ctx         context.Context
		c           client.Client
		reconciler  *ProvisioningRequestReconciler
		task        *provisioningRequestReconcilerTask
		cr          *provisioningv1alpha1.ProvisioningRequest
		ci          *unstructured.Unstructured
		nar         *pluginsv1alpha1.NodeAllocationRequest
		crName      = "cluster-1"
		ctNamespace = "clustertemplate-a-v4-16"
	)

	BeforeEach(func() {
		ctx = context.Background()
		// Define the cluster instance.
		ci = &unstructured.Unstructured{}
		ci.SetName(crName)
		ci.SetNamespace(ctNamespace)
		ci.Object = map[string]interface{}{
			"spec": map[string]interface{}{
				"nodes": []interface{}{
					map[string]interface{}{"role": "master"},
					map[string]interface{}{"role": "master"},
					map[string]interface{}{"role": "worker"},
				},
			},
		}

		// Define the provisioning request.
		cr = &provisioningv1alpha1.ProvisioningRequest{
			ObjectMeta: metav1.ObjectMeta{
				Name: crName,
			},
			Status: provisioningv1alpha1.ProvisioningRequestStatus{
				Extensions: provisioningv1alpha1.Extensions{
					NodeAllocationRequestRef: &provisioningv1alpha1.NodeAllocationRequestRef{
						NodeAllocationRequestID:        crName,
						HardwareProvisioningCheckStart: &metav1.Time{Time: time.Now()},
					},
				},
			},
		}

		// Define the node allocation request.
		nar = &pluginsv1alpha1.NodeAllocationRequest{
			ObjectMeta: metav1.ObjectMeta{
				Name: crName,
			},
			// Set up your NodeAllocationRequest object as needed
			Status: pluginsv1alpha1.NodeAllocationRequestStatus{
				Conditions: []metav1.Condition{},
			},
		}

		c = getFakeClientFromObjects([]client.Object{cr}...)
		reconciler = &ProvisioningRequestReconciler{
			Client: c,
			Logger: logger,
		}
		task = &provisioningRequestReconcilerTask{
			logger: reconciler.Logger,
			client: reconciler.Client,
			object: cr,
			timeouts: &timeouts{
				hardwareProvisioning: 1 * time.Minute,
			},
		}
	})

	It("returns error when error fetching NodeAllocationRequest", func() {
		provisioned, timedOutOrFailed, err := task.checkNodeAllocationRequestStatus(ctx, hwmgmtv1alpha1.Provisioned)
		Expect(provisioned).To(Equal(false))
		Expect(timedOutOrFailed).To(Equal(false))
		Expect(err).To(HaveOccurred())
	})

	It("returns failed when NodeAllocationRequest provisioning failed", func() {
		provisionedCondition := metav1.Condition{
			Type:   "Provisioned",
			Status: metav1.ConditionFalse,
			Reason: string(hwmgmtv1alpha1.Failed),
		}
		nar.Status.Conditions = append(nar.Status.Conditions, provisionedCondition)
		Expect(c.Create(ctx, nar)).To(Succeed())
		provisioned, timedOutOrFailed, err := task.checkNodeAllocationRequestStatus(ctx, hwmgmtv1alpha1.Provisioned)
		Expect(provisioned).To(Equal(false))
		Expect(timedOutOrFailed).To(Equal(true)) // It should be failed
		Expect(err).ToNot(HaveOccurred())
		condition := meta.FindStatusCondition(cr.Status.Conditions, string(provisioningv1alpha1.PRconditionTypes.HardwareProvisioned))
		Expect(condition).ToNot(BeNil())
		Expect(condition.Status).To(Equal(metav1.ConditionFalse))
		Expect(condition.Reason).To(Equal(string(hwmgmtv1alpha1.Failed)))
	})

	It("returns timeout when NodeAllocationRequest provisioning timed out", func() {
		provisionedCondition := metav1.Condition{
			Type:   "Provisioned",
			Status: metav1.ConditionFalse,
		}
		nar.Status.Conditions = append(nar.Status.Conditions, provisionedCondition)
		Expect(c.Create(ctx, nar)).To(Succeed())

		// First call to checkNodeAllocationRequestStatus (before timeout)
		provisioned, timedOutOrFailed, err := task.checkNodeAllocationRequestStatus(ctx, hwmgmtv1alpha1.Provisioned)
		Expect(provisioned).To(Equal(false))
		Expect(timedOutOrFailed).To(Equal(false))
		Expect(err).ToNot(HaveOccurred())

		// Simulate a timeout by moving the start time back
		adjustedTime := cr.Status.Extensions.NodeAllocationRequestRef.HardwareProvisioningCheckStart.Time.Add(-1 * time.Minute)
		cr.Status.Extensions.NodeAllocationRequestRef.HardwareProvisioningCheckStart = &metav1.Time{Time: adjustedTime}

		// Call checkNodeAllocationRequestStatus again (after timeout)
		provisioned, timedOutOrFailed, err = task.checkNodeAllocationRequestStatus(ctx, hwmgmtv1alpha1.Provisioned)
		Expect(provisioned).To(Equal(false))
		Expect(timedOutOrFailed).To(Equal(true)) // Now it should time out
		Expect(err).ToNot(HaveOccurred())

		condition := meta.FindStatusCondition(cr.Status.Conditions, string(provisioningv1alpha1.PRconditionTypes.HardwareProvisioned))
		Expect(condition).ToNot(BeNil())
		Expect(condition.Status).To(Equal(metav1.ConditionFalse))
		Expect(condition.Reason).To(Equal(string(hwmgmtv1alpha1.TimedOut)))
	})

	It("returns false when NodeAllocationRequest is not provisioned", func() {
		provisionedCondition := metav1.Condition{
			Type:   "Provisioned",
			Status: metav1.ConditionFalse,
		}
		nar.Status.Conditions = append(nar.Status.Conditions, provisionedCondition)
		Expect(c.Create(ctx, nar)).To(Succeed())

		provisioned, timedOutOrFailed, err := task.checkNodeAllocationRequestStatus(ctx, hwmgmtv1alpha1.Provisioned)
		Expect(provisioned).To(Equal(false))
		Expect(timedOutOrFailed).To(Equal(false))
		Expect(err).ToNot(HaveOccurred())
		condition := meta.FindStatusCondition(cr.Status.Conditions, string(provisioningv1alpha1.PRconditionTypes.HardwareProvisioned))
		Expect(condition).ToNot(BeNil())
		Expect(condition.Status).To(Equal(metav1.ConditionFalse))
	})

	It("returns true when NodeAllocationRequest is provisioned", func() {
		provisionedCondition := metav1.Condition{
			Type:   "Provisioned",
			Status: metav1.ConditionTrue,
		}
		nar.Status.Conditions = append(nar.Status.Conditions, provisionedCondition)
		Expect(c.Create(ctx, nar)).To(Succeed())
		provisioned, timedOutOrFailed, err := task.checkNodeAllocationRequestStatus(ctx, hwmgmtv1alpha1.Provisioned)
		Expect(provisioned).To(Equal(true))
		Expect(timedOutOrFailed).To(Equal(false))
		Expect(err).ToNot(HaveOccurred())
		condition := meta.FindStatusCondition(cr.Status.Conditions, string(provisioningv1alpha1.PRconditionTypes.HardwareProvisioned))
		Expect(condition).ToNot(BeNil())
		Expect(condition.Status).To(Equal(metav1.ConditionTrue))
	})

	It("returns timeout when NodeAllocationRequest configuring timed out", func() {
		// Set the configuration start time.
		cr.Status.Extensions.NodeAllocationRequestRef.HardwareConfiguringCheckStart = &metav1.Time{Time: time.Now()}
		Expect(c.Status().Update(ctx, cr)).To(Succeed())

		provisionedCondition := metav1.Condition{
			Type:   "Provisioned",
			Status: metav1.ConditionTrue,
		}
		nar.Status.Conditions = append(nar.Status.Conditions, provisionedCondition)
		Expect(c.Create(ctx, nar)).To(Succeed())

		configuredCondition := metav1.Condition{
			Type:   "Configured",
			Status: metav1.ConditionFalse,
		}
		nar.Status.Conditions = append(nar.Status.Conditions, configuredCondition)
		Expect(c.Status().Update(ctx, nar)).To(Succeed())

		// First call to checkNodeAllocationRequestStatus (before timeout)
		status, timedOutOrFailed, err := task.checkNodeAllocationRequestStatus(ctx, hwmgmtv1alpha1.Configured)
		Expect(status).To(Equal(false))
		Expect(timedOutOrFailed).To(Equal(false))
		Expect(err).ToNot(HaveOccurred())

		// Simulate a timeout by moving the start time back
		adjustedTime := cr.Status.Extensions.NodeAllocationRequestRef.HardwareConfiguringCheckStart.Time.Add(-1 * time.Minute)
		cr.Status.Extensions.NodeAllocationRequestRef.HardwareConfiguringCheckStart = &metav1.Time{Time: adjustedTime}

		// Call checkNodeAllocationRequestStatus again (after timeout)
		status, timedOutOrFailed, err = task.checkNodeAllocationRequestStatus(ctx, hwmgmtv1alpha1.Configured)
		Expect(status).To(Equal(false))
		Expect(timedOutOrFailed).To(Equal(true)) // Now it should time out
		Expect(err).ToNot(HaveOccurred())

		condition := meta.FindStatusCondition(cr.Status.Conditions, string(provisioningv1alpha1.PRconditionTypes.HardwareConfigured))
		Expect(condition).ToNot(BeNil())
		Expect(condition.Status).To(Equal(metav1.ConditionFalse))
		Expect(condition.Reason).To(Equal(string(hwmgmtv1alpha1.TimedOut)))
	})
})

var _ = Describe("updateClusterInstance", func() {
	var (
		ctx         context.Context
		c           client.Client
		reconciler  *ProvisioningRequestReconciler
		task        *provisioningRequestReconcilerTask
		cr          *provisioningv1alpha1.ProvisioningRequest
		ci          *siteconfig.ClusterInstance
		nar         *pluginsv1alpha1.NodeAllocationRequest
		tmpNar      *hwmgrpluginapi.NodeAllocationRequestResponse
		crName      = "cluster-1"
		crNamespace = "clustertemplate-a-v4-16"
		mn          = "master-node"
		wn          = "worker-node"
		mhost       = "node1.test.com"
		whost       = "node2.test.com"
		poolns      = utils.UnitTestHwmgrNamespace
		mIfaces     = []*pluginsv1alpha1.Interface{
			{
				Name:       "eth0",
				Label:      "test",
				MACAddress: "00:00:00:01:20:30",
			},
			{
				Name:       "eth1",
				Label:      "test2",
				MACAddress: "66:77:88:99:CC:BB",
			},
		}
		wIfaces = []*pluginsv1alpha1.Interface{
			{
				Name:       "eno1",
				Label:      "test",
				MACAddress: "00:00:00:01:30:10",
			},
			{
				Name:       "eno2",
				Label:      "test2",
				MACAddress: "66:77:88:99:AA:BB",
			},
		}
		masterNode = testutils.CreateNode(mn, "idrac-virtualmedia+https://10.16.2.1/redfish/v1/Systems/System.Embedded.1",
			"site-1-master-bmc-secret", groupNameController, poolns, crName, mIfaces)
		workerNode = testutils.CreateNode(wn, "idrac-virtualmedia+https://10.16.3.4/redfish/v1/Systems/System.Embedded.1",
			"site-1-worker-bmc-secret", groupNameWorker, poolns, crName, wIfaces)
	)

	BeforeEach(func() {
		ctx = context.Background()

		// Define the cluster instance.
		ci = &siteconfig.ClusterInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      crName,
				Namespace: crNamespace,
			},
			Spec: siteconfig.ClusterInstanceSpec{
				Nodes: []siteconfig.NodeSpec{
					{
						Role: "master", HostName: mhost,
						NodeNetwork: &v1beta1.NMStateConfigSpec{
							Interfaces: []*v1beta1.Interface{
								{Name: "eth0"}, {Name: "eth1"},
							},
						},
					},
					{
						Role: "worker", HostName: whost,
						NodeNetwork: &v1beta1.NMStateConfigSpec{
							Interfaces: []*v1beta1.Interface{
								{Name: "eno1"}, {Name: "eno2"},
							},
						},
					},
				},
			},
		}

		// Define the provisioning request.
		cr = &provisioningv1alpha1.ProvisioningRequest{
			ObjectMeta: metav1.ObjectMeta{
				Name: crName,
			},
		}

		// Define the node allocation request.
		nar = &pluginsv1alpha1.NodeAllocationRequest{
			ObjectMeta: metav1.ObjectMeta{
				Name:      crName,
				Namespace: poolns,
				Annotations: map[string]string{
					hwmgmtv1alpha1.BootInterfaceLabelAnnotation: "test",
				},
			},
			Status: pluginsv1alpha1.NodeAllocationRequestStatus{
				Conditions: []metav1.Condition{
					{
						Type:   "Provisioned",
						Status: "True",
					},
				},
				Properties: hwmgmtv1alpha1.Properties{
					NodeNames: []string{mn, wn},
				},
			},
			Spec: pluginsv1alpha1.NodeAllocationRequestSpec{
				NodeGroup: []pluginsv1alpha1.NodeGroup{
					{
						NodeGroupData: hwmgmtv1alpha1.NodeGroupData{
							Name: groupNameController,
							Role: "master",
						},
					},
					{
						NodeGroupData: hwmgmtv1alpha1.NodeGroupData{
							Name: groupNameWorker,
							Role: "worker",
						},
					},
				},
			},
		}

		tmpNar = &hwmgrpluginapi.NodeAllocationRequestResponse{}

		c = getFakeClientFromObjects([]client.Object{cr}...)
		reconciler = &ProvisioningRequestReconciler{
			Client: c,
			Logger: logger,
		}
		task = &provisioningRequestReconcilerTask{
			logger:       reconciler.Logger,
			client:       reconciler.Client,
			object:       cr,
			clusterInput: &clusterInput{},
		}
	})

	It("returns error when failing to get the Node object", func() {
		unstructuredCi, err := utils.ConvertToUnstructured(*ci)
		Expect(err).ToNot(HaveOccurred())
		err = task.updateClusterInstance(ctx, unstructuredCi, tmpNar)
		Expect(err).To(HaveOccurred())
	})

	It("returns error when no match hardware node", func() {
		task.clusterInput.clusterInstanceData = map[string]any{
			"nodes": []any{
				map[string]any{
					"hostName": "masterNode",
				},
				map[string]any{
					"hostName": "workerNode",
				},
			},
		}
		nar.Status.Properties = hwmgmtv1alpha1.Properties{}
		unstructuredCi, err := utils.ConvertToUnstructured(*ci)
		Expect(err).ToNot(HaveOccurred())
		err = task.updateClusterInstance(ctx, unstructuredCi, tmpNar)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("failed to find matches for the following nodes"))
	})

	It("returns no error when updateClusterInstance succeeds", func() {
		task.clusterInput.clusterInstanceData = map[string]any{
			"nodes": []any{
				map[string]any{
					"hostName": mhost,
					"nodeNetwork": map[string]any{
						"interfaces": []any{
							map[string]any{
								"name":  "eth0",
								"label": "test",
							},
							map[string]any{
								"name":  "eth1",
								"label": "test2",
							},
						},
					},
				},
				map[string]any{
					"hostName": whost,
					"nodeNetwork": map[string]any{
						"interfaces": []any{
							map[string]any{
								"name":  "eno1",
								"label": "test",
							},
							map[string]any{
								"name":  "eno2",
								"label": "test2",
							},
						},
					},
				},
			},
		}
		nodes := []*pluginsv1alpha1.AllocatedNode{masterNode, workerNode}
		secrets := testutils.CreateSecrets([]string{masterNode.Status.BMC.CredentialsName, workerNode.Status.BMC.CredentialsName}, poolns)

		testutils.CreateResources(ctx, c, nodes, secrets)
		unstructuredCi, err := utils.ConvertToUnstructured(*ci)
		Expect(err).ToNot(HaveOccurred())
		err = task.updateClusterInstance(ctx, unstructuredCi, tmpNar)
		Expect(err).ToNot(HaveOccurred())

		masterBootMAC, err := utils.GetBootMacAddress(masterNode.Status.Interfaces, "")
		Expect(err).ToNot(HaveOccurred())
		workerBootMAC, err := utils.GetBootMacAddress(workerNode.Status.Interfaces, "")
		Expect(err).ToNot(HaveOccurred())

		// Define expected details
		expectedDetails := []expectedNodeDetails{
			{
				BMCAddress:         masterNode.Status.BMC.Address,
				BMCCredentialsName: masterNode.Status.BMC.CredentialsName,
				BootMACAddress:     masterBootMAC,
				Interfaces:         getInterfaceMap(masterNode.Status.Interfaces),
			},
			{
				BMCAddress:         workerNode.Status.BMC.Address,
				BMCCredentialsName: workerNode.Status.BMC.CredentialsName,
				BootMACAddress:     workerBootMAC,
				Interfaces:         getInterfaceMap(workerNode.Status.Interfaces),
			},
		}

		// Verify the bmc address, secret, boot mac address and interface mac addresses are set correctly in the cluster instance
		ci, err := utils.ConvertFromUnstructured(unstructuredCi)
		Expect(err).ToNot(HaveOccurred())
		verifyClusterInstance(ci, expectedDetails)

		// Verify the host name is set in the node status
		verifyNodeStatus(ctx, c, nodes, mhost, whost)
	})
})

// Helper function to transform interfaces into the required map[string]interface{} format
func getInterfaceMap(interfaces []*pluginsv1alpha1.Interface) []map[string]interface{} {
	var ifaceList []map[string]interface{}
	for _, iface := range interfaces {
		ifaceList = append(ifaceList, map[string]interface{}{
			"Name":       iface.Name,
			"MACAddress": iface.MACAddress,
		})
	}
	return ifaceList
}


func verifyClusterInstance(ci *siteconfig.ClusterInstance, expectedDetails []expectedNodeDetails) {
	for i, expected := range expectedDetails {
		Expect(ci.Spec.Nodes[i].BmcAddress).To(Equal(expected.BMCAddress))
		Expect(ci.Spec.Nodes[i].BmcCredentialsName.Name).To(Equal(expected.BMCCredentialsName))
		Expect(ci.Spec.Nodes[i].BootMACAddress).To(Equal(expected.BootMACAddress))
		// Verify Interface MAC Address
		for _, iface := range ci.Spec.Nodes[i].NodeNetwork.Interfaces {
			for _, expectedIface := range expected.Interfaces {
				// Compare the interface name and MAC address
				if iface.Name == expectedIface["Name"] {
					Expect(iface.MacAddress).To(Equal(expectedIface["MACAddress"]), "MAC Address mismatch for interface")
				}
			}
		}
	}
}

func verifyNodeStatus(ctx context.Context, c client.Client, nodes []*pluginsv1alpha1.AllocatedNode, mhost, whost string) {
	for _, node := range nodes {
		updatedNode := &pluginsv1alpha1.AllocatedNode{}
		Expect(c.Get(ctx, client.ObjectKey{Name: node.Name, Namespace: node.Namespace}, updatedNode)).To(Succeed())
		switch updatedNode.Spec.GroupName {
		case groupNameController:
			Expect(updatedNode.Status.Hostname).To(Equal(mhost))
		case groupNameWorker:
			Expect(updatedNode.Status.Hostname).To(Equal(whost))
		default:
			Fail(fmt.Sprintf("Unexpected GroupName: %s", updatedNode.Spec.GroupName))
		}
	}
}
//*/

func VerifyHardwareTemplateStatus(ctx context.Context, c client.Client, templateName string, expectedCon metav1.Condition) {
	updatedHwTempl := &hwmgmtv1alpha1.HardwareTemplate{}
	Expect(c.Get(ctx, client.ObjectKey{Name: templateName, Namespace: utils.InventoryNamespace}, updatedHwTempl)).To(Succeed())
	hwTemplCond := meta.FindStatusCondition(updatedHwTempl.Status.Conditions, expectedCon.Type)
	Expect(hwTemplCond).ToNot(BeNil())
	testutils.VerifyStatusCondition(*hwTemplCond, metav1.Condition{
		Type:    expectedCon.Type,
		Status:  expectedCon.Status,
		Reason:  expectedCon.Reason,
		Message: expectedCon.Message,
	})
}
