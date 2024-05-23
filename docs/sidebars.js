/**
 * Creating a sidebar enables you to:
 - create an ordered group of docs
 - render a sidebar for each doc of that group
 - provide next/previous navigation

 The sidebars can be generated from the filesystem, or explicitly defined here.

 Create as many sidebars as you want.
 */

// @ts-check

/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  // By default, Docusaurus generates a sidebar from the docs folder structure
  // tutorialSidebar: [{type: 'autogenerated', dirName: '.'}],

  // But you can create a sidebar manually
  docs: [
    {
      type: 'doc',
      label: 'Introduction',
      id: 'intro'
    },
    {
      type: 'category',
      label: 'Getting started',
      link: {
        type: 'generated-index',
      },
      items: [
        {
          type: 'doc',
          label: 'Quickstart',
          id: 'getting-started/quickstart',
        },
        {
          type: 'doc',
          label: 'Quickstart SGX',
          id: 'getting-started/quickstart-sgx',
        },
        {
          type: 'doc',
          label: 'Quickstart Simulation',
          id: 'getting-started/quickstart-simulation',
        },
        {
          type: 'doc',
          label: 'Concepts',
          id: 'getting-started/concepts',
        },
        {
          type: 'doc',
          label: 'Coordinator',
          id: 'getting-started/coordinator',
        },
        {
          type: 'doc',
          label: 'Marbles',
          id: 'getting-started/marbles',
        },
      ],
    },
    {
      type: 'category',
      label: 'Features',
      link: {
        type: 'generated-index',
      },
      items: [
        {
          type: 'doc',
          label: 'Manifest',
          id: 'features/manifest',
        },
        {
          type: 'doc',
          label: 'Attestation',
          id: 'features/attestation',
        },
        {
          type: 'doc',
          label: 'State and recovery',
          id: 'features/recovery',
        },
        {
          type: 'doc',
          label: 'Secrets management',
          id: 'features/secrets-management',
        },
        {
          type: 'doc',
          label: 'Transparent TLS',
          id: 'features/transparent-TLS',
        },
        {
          type: 'doc',
          label: 'Kubernetes integration',
          id: 'features/kubernetes-integration',
        },
        {
          type: 'doc',
          label: 'Supported runtimes',
          id: 'features/runtimes',
        },
      ],
    },
    {
      type: 'category',
      label: 'Deployment',
      link: {
        type: 'generated-index',
      },
      items: [
        {
          type: 'doc',
          label: 'Cloud',
          id: 'deployment/cloud',
        },
        {
          type: 'doc',
          label: 'On-premises',
          id: 'deployment/on-prem',
        },
        {
          type: 'doc',
          label: 'Kubernetes',
          id: 'deployment/kubernetes',
        },
        {
          type: 'doc',
          label: 'Standalone',
          id: 'deployment/standalone',
        },
      ],
    }, {
      type: 'category',
      label: 'Workflows',
      link: {
        type: 'generated-index',
      },
      items: [
        {
          type: 'doc',
          label: 'Define a manifest',
          id: 'workflows/define-manifest',
        },
        {
          type: 'doc',
          label: 'Set a manifest',
          id: 'workflows/set-manifest',
        },
        {
          type: 'doc',
          label: 'Add a service',
          id: 'workflows/add-service',
        },
        {
          type: 'doc',
          label: 'Verify a deployment',
          id: 'workflows/verification',
        },
        {
          type: 'doc',
          label: 'Monitoring and logging',
          id: 'workflows/monitoring',
        },
        {
          type: 'doc',
          label: 'Update a manifest',
          id: 'workflows/update-manifest',
        },
        {
          type: 'doc',
          label: 'Manage secrets',
          id: 'workflows/managing-secrets',
        },
        {
          type: 'doc',
          label: 'Update a deployment',
          id: 'workflows/updates',
        },
        {
          type: 'doc',
          label: 'Recover the Coordinator',
          id: 'workflows/recover-coordinator',
        },
      ],
    },
    {
      type: 'category',
      label: 'Building services',
      link: {
        type: 'generated-index',
      },
      items: [
        {
          type: 'doc',
          label: 'EGo',
          id: 'building-services/ego',
        },
        {
          type: 'doc',
          label: 'Gramine',
          id: 'building-services/gramine',
        },
        {
          type: 'doc',
          label: 'Occlum',
          id: 'building-services/occlum',
        },
      ],
    },
    {
      type: 'doc',
      label: 'Examples',
      id: 'examples',
    },
    {
      type: 'category',
      label: 'Reference',
      link: {
        type: 'generated-index',
      },
      items: [
        {
          type: 'doc',
          label: 'Coordinator client API',
          id: 'reference/coordinator',
        },
        {
          type: 'doc',
          label: 'CLI',
          id: 'reference/cli',
        },
      ],
    },
  ],
};

module.exports = sidebars;
