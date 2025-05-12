# [0.20.0](https://github.com/nestoca/joy-generator/compare/v0.19.2...v0.20.0) (2025-05-12)


### Features

* **PL-3383:** add trace specifically for ComputeReleaseValues ([9598b11](https://github.com/nestoca/joy-generator/commit/9598b11ae0963100718bee6f4f6b7c0c0728fd92))



## [0.19.2](https://github.com/nestoca/joy-generator/compare/v0.19.1...v0.19.2) (2025-04-09)


### Bug Fixes

* **ci:** update deploy pages artifact step ([b4920d7](https://github.com/nestoca/joy-generator/commit/b4920d720c56ac2f36431d288ddd04755700d73e))



## [0.19.1](https://github.com/nestoca/joy-generator/compare/v0.19.0...v0.19.1) (2025-04-09)


### Bug Fixes

* use joy-ci-actions github app ([49c85f8](https://github.com/nestoca/joy-generator/commit/49c85f880bd2f1f926184e0a139c4445042044a8))



# [0.19.0](https://github.com/nestoca/joy-generator/compare/v0.18.1...v0.19.0) (2025-04-04)


### Bug Fixes

* **PL-3219:** update (deprecated) actions/upload-pages-artifact to v4 ([#73](https://github.com/nestoca/joy-generator/issues/73)) ([8d3d41c](https://github.com/nestoca/joy-generator/commit/8d3d41c97c7d3d68e267b45cdf262003c30bc1c1))
* **PL-3291:** update direct/indirect vulnerable dependencies ([#71](https://github.com/nestoca/joy-generator/issues/71)) ([a2ce366](https://github.com/nestoca/joy-generator/commit/a2ce3669cad6023e2dd909c9bc3dd29f8f199de3))


### Features

* **PL-3378:** remove secrets.GH_TOKEN refs ([96367d3](https://github.com/nestoca/joy-generator/commit/96367d35d24815272d48f7bef8fdc789bcfed332))



## [0.18.1](https://github.com/nestoca/joy-generator/compare/v0.18.0...v0.18.1) (2024-11-12)


### Bug Fixes

* **PL-3111:** use joy 0.63.1 ([#70](https://github.com/nestoca/joy-generator/issues/70)) ([92cd284](https://github.com/nestoca/joy-generator/commit/92cd284eec986d519c5c31b2a915d22dbb7c0d26))



# [0.18.0](https://github.com/nestoca/joy-generator/compare/v0.17.0...v0.18.0) (2024-10-30)


### Features

* **PL-3082:** add namespace field to releases ([#68](https://github.com/nestoca/joy-generator/issues/68)) ([244d3b9](https://github.com/nestoca/joy-generator/commit/244d3b96a8af77080acc1c54deefe92d21dbc68d))



# [0.17.0](https://github.com/nestoca/joy-generator/compare/v0.16.1...v0.17.0) (2024-10-30)


### Features

* **PL-3082:** add namespace field to releases ([#67](https://github.com/nestoca/joy-generator/issues/67)) ([9b9850f](https://github.com/nestoca/joy-generator/commit/9b9850f7c88f44b725084b8ce6226ee5538b925d))



## [0.16.1](https://github.com/nestoca/joy-generator/compare/v0.16.0...v0.16.1) (2024-10-10)


### Bug Fixes

* **deps:** update module github.com/nestoca/joy to v0.61.1 ([d4ddb29](https://github.com/nestoca/joy-generator/commit/d4ddb29f0e3a2bf8a36d82996e59d04bcb2587b3))



# [0.16.0](https://github.com/nestoca/joy-generator/compare/v0.15.0...v0.16.0) (2024-09-03)


### Features

* **PL-2911:** configure health probe timeouts ([28bd67d](https://github.com/nestoca/joy-generator/commit/28bd67df05dc2ae99280f8d93e498cabef048d4b))



# [0.15.0](https://github.com/nestoca/joy-generator/compare/v0.14.0...v0.15.0) (2024-08-30)


### Features

* release-version ([33dd66b](https://github.com/nestoca/joy-generator/commit/33dd66bf43e5093ab05b44bd14f4d477dfab8f9c))



# [0.14.0](https://github.com/nestoca/joy-generator/compare/v0.13.0...v0.14.0) (2024-08-30)


### Bug Fixes

* minimum concurrency ([6c7a750](https://github.com/nestoca/joy-generator/commit/6c7a7502d2e8f55f7a990bb66131c97fdffa7900))
* span names ([5b6ee95](https://github.com/nestoca/joy-generator/commit/5b6ee95def112084f79d7ed5d51a7825e78b7f0e))


### Features

* make generator concurrency configurable ([92a3ab8](https://github.com/nestoca/joy-generator/commit/92a3ab8e05f80856d7d62c66dc5c0af5b32443e0))



# [0.13.0](https://github.com/nestoca/joy-generator/compare/v0.12.0...v0.13.0) (2024-08-30)


### Features

* add concurrency to release rendering ([c364039](https://github.com/nestoca/joy-generator/commit/c3640397ab9f31ff7713c6981f7cdb3a2cbbc310))



# [0.12.0](https://github.com/nestoca/joy-generator/compare/v0.11.0...v0.12.0) (2024-08-29)


### Bug Fixes

* **PL-2800:** downgrade otelhttp because of introduced bug superflously writing header ([507f45a](https://github.com/nestoca/joy-generator/commit/507f45a720b9b6d2a22db9355faedd703a2af7a4))


### Features

* **PL-2800:** add spans to generator run ([218b1c3](https://github.com/nestoca/joy-generator/commit/218b1c3eceb3e57c6306847b95934269d5905899))



# [0.11.0](https://github.com/nestoca/joy-generator/compare/v0.10.0...v0.11.0) (2024-08-29)


### Features

* **PL-2800:** setup opentelemetry tracing ([ef190e4](https://github.com/nestoca/joy-generator/commit/ef190e4857cfad70888062699084764030a5b344))



# [0.10.0](https://github.com/nestoca/joy-generator/compare/v0.9.0...v0.10.0) (2024-07-25)


### Features

* **PL-2834:** trigger publish pipeline ([1cf1c8b](https://github.com/nestoca/joy-generator/commit/1cf1c8b8996d01796484c8aaae725393ec7d91a4))
* **PL-2834:** trigger publish pipeline ([ac11623](https://github.com/nestoca/joy-generator/commit/ac1162373c11e5a2774a132e659c1817dd471195))



# [0.9.0](https://github.com/nestoca/joy-generator/compare/v0.8.5...v0.9.0) (2024-07-08)


### Features

* **PL-2806:** make joy-generator requestTimeout configurable ([0c0d5fe](https://github.com/nestoca/joy-generator/commit/0c0d5fed6838306baffc1ce741a6d32b75a1761a))



## [0.8.5](https://github.com/nestoca/joy-generator/compare/v0.8.4...v0.8.5) (2024-06-19)


### Bug Fixes

* **deps:** update module github.com/nestoca/joy to v0.55.1 ([58e1819](https://github.com/nestoca/joy-generator/commit/58e181950274884e5b87d1e049de69b4b4018c41))



## [0.8.4](https://github.com/nestoca/joy-generator/compare/v0.8.3...v0.8.4) (2024-06-18)


### Bug Fixes

* **PL-2763:** check cache folder is empty before pulling ([9fe62df](https://github.com/nestoca/joy-generator/commit/9fe62dffccf6ac571dfef2f85161f59478cc8fb9))



## [0.8.3](https://github.com/nestoca/joy-generator/compare/v0.8.2...v0.8.3) (2024-06-18)


### Bug Fixes

* **PL-2763:** check chart has been written before deciding to pull ([19d5fa3](https://github.com/nestoca/joy-generator/commit/19d5fa3b1aa2834c5f166623aa4a044671ad5528))



## [0.8.2](https://github.com/nestoca/joy-generator/compare/v0.8.1...v0.8.2) (2024-06-12)


### Bug Fixes

* **PL-2763:** do not pull charts concurrently ([707f547](https://github.com/nestoca/joy-generator/commit/707f547c9ed56b3f8daa115806dc1b9e1641313d))



## [0.8.1](https://github.com/nestoca/joy-generator/compare/v0.8.0...v0.8.1) (2024-06-04)


### Bug Fixes

* **deps:** update module github.com/nestoca/joy to v0.54.5 ([21c25b6](https://github.com/nestoca/joy-generator/commit/21c25b6acdb479f0b74dbf63d2529bbc169ce086))



# [0.8.0](https://github.com/nestoca/joy-generator/compare/v0.7.0...v0.8.0) (2024-05-24)


### Features

* **PL-2701:** revert "update joy to unify schema values" ([3d2480c](https://github.com/nestoca/joy-generator/commit/3d2480cc985a43d23b77afa4061256728633465a))



# [0.7.0](https://github.com/nestoca/joy-generator/compare/v0.6.3...v0.7.0) (2024-05-22)


### Features

* **PL-2701:** Revert "update joy to unify schema values" ([3cc8621](https://github.com/nestoca/joy-generator/commit/3cc86215d3380222b5cce553f96ca87646cca980))
* **PL-2729:** Update joy to support .joyignore ([5ce3282](https://github.com/nestoca/joy-generator/commit/5ce32825a67d0168bf640942e716460161b8cef4))



## [0.6.3](https://github.com/nestoca/joy-generator/compare/v0.6.2...v0.6.3) (2024-05-13)


### Bug Fixes

* **PL-2701:** fix credential read operation ([8931ab2](https://github.com/nestoca/joy-generator/commit/8931ab2459a8eef431a387acff78052c9ae34b24))



## [0.6.2](https://github.com/nestoca/joy-generator/compare/v0.6.1...v0.6.2) (2024-05-13)


### Bug Fixes

* **PL-2701:** log helm authentication success event ([e559d23](https://github.com/nestoca/joy-generator/commit/e559d238a4a44d803117cabde3eeb6d0a5e07834))



## [0.6.1](https://github.com/nestoca/joy-generator/compare/v0.6.0...v0.6.1) (2024-05-13)


### Bug Fixes

* **PL-2701:** give user permission to home for storing helm cache ([31917ec](https://github.com/nestoca/joy-generator/commit/31917ec85be1864fc9826805068711c1b22ca64f))



# [0.6.0](https://github.com/nestoca/joy-generator/compare/v0.5.12...v0.6.0) (2024-05-13)


### Features

* **PL-2701:** update joy to unify schema values ([156132f](https://github.com/nestoca/joy-generator/commit/156132fbcc48abb5bc4b8ad038e9b98b65272695))



## [0.5.12](https://github.com/nestoca/joy-generator/compare/v0.5.11...v0.5.12) (2024-05-10)


### Bug Fixes

* add datadog service label value to chart ([fe50929](https://github.com/nestoca/joy-generator/commit/fe5092973a2c5f38b965c8cebc513bd1f14a8641))



## [0.5.11](https://github.com/nestoca/joy-generator/compare/v0.5.10...v0.5.11) (2024-05-07)


### Bug Fixes

* **deps:** update module github.com/nestoca/joy to v0.50.0 ([#40](https://github.com/nestoca/joy-generator/issues/40)) ([40a61a3](https://github.com/nestoca/joy-generator/commit/40a61a3e600692fbc9d617f74b3df88d0df27d72))



## [0.5.10](https://github.com/nestoca/joy-generator/compare/v0.5.9...v0.5.10) (2024-05-02)


### Bug Fixes

* **deps:** update module github.com/nestoca/joy to v0.49.3 ([96a1800](https://github.com/nestoca/joy-generator/commit/96a18005686dd11d3e5133bf6f6d06a42c521f66))



## [0.5.9](https://github.com/nestoca/joy-generator/compare/v0.5.8...v0.5.9) (2024-04-22)


### Bug Fixes

* **PL-2654:** stop logging context canceled on shutdown to stderr ([219295c](https://github.com/nestoca/joy-generator/commit/219295c3b513da4422357462fbabd3fef7b91a80))



## [0.5.8](https://github.com/nestoca/joy-generator/compare/v0.5.7...v0.5.8) (2024-04-10)


### Bug Fixes

* **deps:** update module github.com/nestoca/joy to v0.47.5 ([82f6e01](https://github.com/nestoca/joy-generator/commit/82f6e01afa68d09ffe99fea773474ad30a3018be))



## [0.5.7](https://github.com/nestoca/joy-generator/compare/v0.5.6...v0.5.7) (2024-04-05)


### Bug Fixes

* **PL-2579:** renovate pr for joy only ([#20](https://github.com/nestoca/joy-generator/issues/20)) ([61c9d18](https://github.com/nestoca/joy-generator/commit/61c9d18626159be34998de0c22543c54662cdd75))



## [0.5.6](https://github.com/nestoca/joy-generator/compare/v0.5.5...v0.5.6) (2024-03-25)


### Bug Fixes

* **PL-2554:** fix git reset logic ([4323625](https://github.com/nestoca/joy-generator/commit/4323625c0e8469288a4bc63d79e27f8b2d601683))



## [0.5.5](https://github.com/nestoca/joy-generator/compare/v0.5.4...v0.5.5) (2024-03-25)


### Bug Fixes

* **PL-2554:** capture stack trace for recovery panics ([680db09](https://github.com/nestoca/joy-generator/commit/680db0910a29bcf1f93ce206b96aed0e584a146a))



## [0.5.4](https://github.com/nestoca/joy-generator/compare/v0.5.3...v0.5.4) (2024-03-25)


### Bug Fixes

* **PL-2554:** fix temp test issue ([22d14ac](https://github.com/nestoca/joy-generator/commit/22d14ac3c35197b19f34819507dcbfd76eab18f0))
* **PL-2554:** update joy version ([c6728e5](https://github.com/nestoca/joy-generator/commit/c6728e56818662df32572233025490ee2e3f6881))



## [0.5.3](https://github.com/nestoca/joy-generator/compare/v0.5.2...v0.5.3) (2024-03-22)


### Bug Fixes

* **PL-2554:** fix key path variable name ([bf45bd9](https://github.com/nestoca/joy-generator/commit/bf45bd9739c69177c5826cb6ba9e09615d72b48b))



## [0.5.2](https://github.com/nestoca/joy-generator/compare/v0.5.1...v0.5.2) (2024-03-22)


### Bug Fixes

* **PL-2554:** separate secrets from envvars in chart ([bc68e58](https://github.com/nestoca/joy-generator/commit/bc68e58ab044c6ccc2dc38a2350c072882a2f650))



## [0.5.1](https://github.com/nestoca/joy-generator/compare/v0.5.0...v0.5.1) (2024-03-21)


### Bug Fixes

* **PL-2554:** fix chart default values ([af253fa](https://github.com/nestoca/joy-generator/commit/af253faf8f27f83d1ec4ed7686b0211c1982a7b0))



# [0.5.0](https://github.com/nestoca/joy-generator/compare/v0.4.0...v0.5.0) (2024-03-21)


### Features

* **PL-2554:** add e2e test ([9eb1906](https://github.com/nestoca/joy-generator/commit/9eb1906ea3fa5634ca360586c89070f8a2cbcf52))
* **PL-2554:** major refactoring ([a1fd069](https://github.com/nestoca/joy-generator/commit/a1fd069acd699672cd15c0f4773c996bf5c6e572))
* **PL-2554:** update chart ([a36a8a1](https://github.com/nestoca/joy-generator/commit/a36a8a1935ce86db98bc53b1b8df0406747dfa08))



# [0.4.0](https://github.com/nestoca/joy-generator/compare/v0.3.1...v0.4.0) (2024-03-19)


### Features

* **PL-2554:** support chart refs ([eac3533](https://github.com/nestoca/joy-generator/commit/eac3533799e0825d6cea886cf5f809a393e28f46))



## [0.3.1](https://github.com/nestoca/joy-generator/compare/v0.3.0...v0.3.1) (2024-03-13)


### Bug Fixes

* **PL-2418:** update joy to support !local tags ([e697712](https://github.com/nestoca/joy-generator/commit/e6977126ca47e08fcbf44ab9df25da60f250fdfa))



# [0.3.0](https://github.com/nestoca/joy-generator/compare/v0.2.2...v0.3.0) (2024-02-22)


### Features

* **PL-2386:** Include project information in results ([e8269c3](https://github.com/nestoca/joy-generator/commit/e8269c3db2bce0065181fbaafa8f459c6027e08c))



## [0.2.2](https://github.com/nestoca/joy-generator/compare/v0.2.1...v0.2.2) (2024-02-20)


### Bug Fixes

* **PL-2372:** update-joy ([111d955](https://github.com/nestoca/joy-generator/commit/111d95557774e84c45fab6a2ebdc984e57221f8e))



## [0.2.1](https://github.com/nestoca/joy-generator/compare/v0.2.0...v0.2.1) (2024-02-19)


### Bug Fixes

* **PL-2372): Revert "fix(PL-2372:** use default gh token" ([e633643](https://github.com/nestoca/joy-generator/commit/e633643f112fb7e0d79ef3a54a13b7e948c2c950))



# [0.1.0](https://github.com/nestoca/joy-generator/compare/v0.0.7...v0.1.0) (2023-11-22)


### Features

* **PL-2009:** Render values with 2-space indents ([b5a119f](https://github.com/nestoca/joy-generator/commit/b5a119f1d15c027436a06920907c989ac173e50c))
* **PL-2009:** Update github.com/nestoca/joy to latest version ([07aa349](https://github.com/nestoca/joy-generator/commit/07aa349c5033268bb67a0bdc4fddcb69d8cc80e0))



## [0.0.7](https://github.com/nestoca/joy-generator/compare/v0.0.6...v0.0.7) (2023-08-21)


### Bug Fixes

* **DEVOPS-1803:** Add annotations for sealedsecrets ([40c25f9](https://github.com/nestoca/joy-generator/commit/40c25f9c75f31d4ecfb499ca0cada9f68a38db14))



## [0.0.6](https://github.com/nestoca/joy-generator/compare/v0.0.5...v0.0.6) (2023-08-18)


### Bug Fixes

* **DEVOPS-1803:** Add readiness probe ([9f5097b](https://github.com/nestoca/joy-generator/commit/9f5097bae20f63b849cf5fe2a31dae086393d3a1))
* **DEVOPS-1803:** Fix bug when remote is force-pushed ([9c94d9d](https://github.com/nestoca/joy-generator/commit/9c94d9daf364e868167a8237905d3b083ca16d4e))
* **DEVOPS-1803:** Improve gitrepo ([b7997bf](https://github.com/nestoca/joy-generator/commit/b7997bf36440c480f1a39d279f3bf889e631da92))
* **DEVOPS-1803:** Set default log level to INFO and add debug flag ([17b96dc](https://github.com/nestoca/joy-generator/commit/17b96dc381d43e58dc556a653dc205012d93a2cc))
* **DEVOPS-1803:** Set GIN_MODE to release to suppress debug logs ([a9a1d24](https://github.com/nestoca/joy-generator/commit/a9a1d24fe7a89476be8b46b818d2056a2a27f480))



## [0.0.5](https://github.com/nestoca/joy-generator/compare/v0.0.4...v0.0.5) (2023-08-17)


### Bug Fixes

* **DEVOPS-1803:** Fix status check & use json logging ([585e68b](https://github.com/nestoca/joy-generator/commit/585e68b7cef1dff84dc03baca467ddd90fbbb1f3))



## [0.0.4](https://github.com/nestoca/joy-generator/compare/v0.0.3...v0.0.4) (2023-08-16)


### Bug Fixes

* **DEVOPS-1803:** Improve config variable loading in chart ([9673288](https://github.com/nestoca/joy-generator/commit/9673288224a048242cbd932fde6561f5d61400cf))



## [0.0.3](https://github.com/nestoca/joy-generator/compare/v0.0.2...v0.0.3) (2023-08-16)


### Bug Fixes

* **DEVOPS-1803:** Various bug fixes ([0cccc01](https://github.com/nestoca/joy-generator/commit/0cccc01ec18556a21591d64d60eb7d1a5d6629ed))



## [0.0.2](https://github.com/nestoca/joy-generator/compare/v0.0.1...v0.0.2) (2023-08-16)


### Bug Fixes

* **DEVOPS-1803:** Fix conditional github app config ([1ed5068](https://github.com/nestoca/joy-generator/commit/1ed5068134699b92818dc47f2b1bcd5cc6485920))



## [0.0.1](https://github.com/nestoca/joy-generator/compare/v0.0.0...v0.0.1) (2023-08-16)


### Bug Fixes

* **DEVOPS-1803:** Fix rendering chart with github token ([6726189](https://github.com/nestoca/joy-generator/commit/67261897618d66227bda41fcc355a5b59f14e4ef))



# 0.0.0 (2023-08-11)


### Features

* **DEVOPS-1803:** Create applicationset generator plugin ([#1](https://github.com/nestoca/joy-generator/issues/1)) ([a63ee63](https://github.com/nestoca/joy-generator/commit/a63ee633b54d3a4399f4b7e23a790f7734250548))
* **DEVOPS-1803:** Publish joy-generator image and chart ([#3](https://github.com/nestoca/joy-generator/issues/3)) ([491b019](https://github.com/nestoca/joy-generator/commit/491b019dc4fb1dea8dd88204698df86d1d64a5ac))



